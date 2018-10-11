# Test Shell-Program for selpg-go
# 1. \n mode, normal file input and stdin input
echo "Case 1=================================================="
selpg-go -s1 -e1 ../test.txt
echo "Case 2=================================================="
selpg-go -s1 -e1 < ../test.txt

# 2. \n mode, for normally, startPage oversize and endPage oversize
echo "Case 3=================================================="
selpg-go -s1 -e2 -l1 ../test.txt
echo "Case 4=================================================="
selpg-go -s4 -e5 -l1 ../test.txt
echo "Case 5=================================================="
selpg-go -s1 -e4 -l1 ../test.txt

# 3. \f mode, also three test cases
echo "Case 6=================================================="
selpg-go -s1 -e2 -f ../testf.txt
echo "Case 7=================================================="
selpg-go -s5 -e6 -f ../testf.txt
echo "Case 8=================================================="
selpg-go -s1 -e5 -f ../testf.txt

# 4. Printer mode, not exist
echo "Case 9=================================================="
selpg-go -s1 -e2 -l1 -dpl ../test.txt