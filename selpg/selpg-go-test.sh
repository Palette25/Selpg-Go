# Test Shell-Program for selpg-go
# 1. \n mode, for normally, startPage oversize and endPage oversize
echo "Case 1=================================================="
selpg-go -s1 -e2 -l1 ../test.txt
echo "Case 2=================================================="
selpg-go -s4 -e5 -l1 ../test.txt
echo "Case 3=================================================="
selpg-go -s1 -e4 -l1 ../test.txt

# 2. \f mode, also three test cases
echo "Case 4=================================================="
selpg-go -s1 -e2 -f ../testf.txt
echo "Case 5=================================================="
selpg-go -s5 -e6 -f ../testf.txt
echo "Case 6=================================================="
selpg-go -s1 -e5 -f ../testf.txt

# 3. Printer mode, not exist
echo "Case 7=================================================="
selpg-go -s1 -e2 -l1 -dpl ../test.txt