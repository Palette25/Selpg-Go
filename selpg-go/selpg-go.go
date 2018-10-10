/*=================================================================

Program name:
	selpg (SELect PaGes)  --- Go Version

Purpose:
	Sometimes one needs to extract only a specified range of
	pages from an input text file. This program allows the user to do
	that.

Author: 
	Palette25
	
Reference: 
	Vasudev Ram -- selpg in C

===================================================================*/

package selpg;

import (
	"fmt"
	"bufio"
	"os"
	"io"
	"os/exec"
	"strings"
	"math"
	flag "github.com/spf13/pflag"
);

/*============================= const values =======================*/
const BUFSIZ = 8192;
const INBUFSIZ = 16 * 1024;
const INT_MAX = int(^uint(0) >> 1);

/*================================= globals =======================*/

var programName string; /* program name, for error messages */

/*========================= structure definition ===================*/
/* Selpg-Args */
type sp_args struct {
	start_page int
	end_page int
	file_name string // uint8(char) array to store file name
	page_len int      // Used for "lines-delimited" command, also as default value
	page_type bool    /* Determine type of page seperator, 'l' for lines-delimited, 
							'f' for form-feed-delimited */
	print_dest string // Destination name
};

func main() {
	// Initialize args and programName
	var (
		startPage, endPage, pageLen int
		pageType, helpFlag bool
		printDest string
	)
	// Binding pflag to params
	flag.IntVarP(&startPage, "start_page", "s", -1, "Starting page of selcetion");
	flag.IntVarP(&endPage, "end_page", "e", -1, "Ending page of selection");
	flag.IntVarP(&pageLen, "page_len", "l", 72, "One page's length");
	flag.BoolVarP(&pageType, "page_type", "f", false, "Using format-limit or not");
	flag.BoolVarP(&helpFlag, "help_flag", "h", false, "Need help");
	flag.StringVarP(&printDest, "print_dest", "d", "", "The Destination file to print");
	flag.Parse();

	programName = os.Args[0];
	/* Check -sstartPage and -eendPage params */
	if startPage == -1 {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page", programName);
		usage();
		os.Exit(8);
	} else if endPage == -1 {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page", programName);
		usage();
		os.Exit(9);
	}

	sp := sp_args{start_page:startPage, end_page:endPage, page_len:pageLen, page_type:pageType, print_dest:printDest};

	process_args(flag.Args(), &sp);
	process_file(sp, helpFlag);

}

/*================================= process_args() ================*/
func process_args(av []string, sp *sp_args) {
	/* Judge startPage and endPage's validness */
	if sp.start_page < 1 || sp.start_page > INT_MAX - 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", programName, sp.start_page);
		usage();
		os.Exit(1);
	}
	if sp.end_page < 1 || sp.end_page > INT_MAX - 1 || sp.end_page < sp.start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", programName, sp.end_page);
		usage();
		os.Exit(2);
	}
	/* Judge pageLength's validness */
	if sp.page_len < 1 || sp.page_len > INT_MAX - 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", programName, sp.page_len);
		usage();
		os.Exit(3);
	}

	/* Get source filename */
	if len(av) == 0 {
		fmt.Fprintf(os.Stderr, "%s: you should enter the source file name", programName);
		usage();
		os.Exit(4);
	}
	tempFileName := av[0];
	_, err := os.Stat(tempFileName);
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: invalid file name -- input file \"%s\" does not exist!\n", programName, tempFileName);
		usage();
		os.Exit(5);
	}
	sp.file_name = tempFileName;
}

/*================================= process_file() ================*/
func process_file(sp sp_args, flag_ bool) {
	if flag_ {
		flag.Usage();
	}
	file, err := os.Open(sp.file_name);
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error with file \"%s\" opening!", programName, sp.file_name);
		os.Exit(6);
	}
	defer file.Close();

	var (
		store, result []string
		targetDelim byte
		leng int
	)
	reader := bufio.NewReader(file);
	pageLength := 0;
	if sp.page_type == true {
		/* Command page type as \f delim a page */
		targetDelim = '\f';
		leng = 1;
	} else {
		targetDelim = '\n';
		leng = sp.page_len;
	}
	/* Loop to get target pages */
	for pageLength < sp.end_page * leng {
		str, err := reader.ReadString(targetDelim);
		if err != nil {
			break;
		}
		pageLength++;
		store = append(store, str);
	}
	/* Check whether the startPage is greater than total Page num */
	if pageLength <= (sp.start_page-1) * leng {
		fmt.Fprintf(os.Stderr, "%s: StartPage(%d) is greater than total pages num(%.0f), no output...\n", programName, sp.start_page, math.Ceil(float64(pageLength)/float64(leng)));
		fmt.Printf("%s: Done\n", programName);
		os.Exit(7);
	}

	/* Getting from startPage to endPage string slice */
	var (
		tempEnd int = sp.end_page * sp.page_len
		tempLen int = sp.page_len
		targetBuf string
	)
	if sp.page_type {
		tempEnd = sp.end_page;
		tempLen = 1;
	}
	if pageLength < sp.end_page * leng {
		fmt.Fprintf(os.Stderr, "%s: EndPage(%d) is greater than total pages num(%.0f), less output than expect....\n", programName, sp.end_page, math.Ceil(float64(pageLength)/float64(leng)));
		tempEnd = len(store);
	}
	result = store[(sp.start_page-1)*tempLen : tempEnd];

	for i:=0; i<len(result); i++ {
		fmt.Printf(result[i]);
		targetBuf = fmt.Sprintf(targetBuf, result[i]);
	}

	/* Printer Destination Output */
	if sp.print_dest != "" {
		cmd := exec.Command("lp", fmt.Sprintf("-d%v", sp.print_dest));
		piper, pipew := io.Pipe();
		stderr,_ := cmd.StderrPipe();
		go func() {
			defer pipew.Close();
			io.Copy(pipew, stderr);
		}()

		cmd.Stdin = strings.NewReader(targetBuf);
		cmd.Run();
		io.Copy(os.Stderr, piper);
		defer bufio.NewWriter(os.Stderr).Flush();
		cmd.Wait();
	}

	fmt.Fprintf(os.Stderr, "%s: Done\n", programName);
}

/*================================= usage() =======================*/
func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", programName);
}