#!/bin/bash

assert() {
    expected="$1"
    input="$2"

    ./main "$input" > tmp.s
    gcc -static -o tmp tmp.s
    ./tmp
    actual="$?"

    if [ "$actual" = "$expected" ]; then
        echo "$input => $actual"
    else
        echo "$input => $expected expected, but got $actual"
        exit 1
    fi
}

assert 0 "0;"
assert 42 "42;"
assert 21 "5+20-4;"
assert 41 " 12 + 34 - 5;"
assert 47 "5+6*7;"
assert 15 "5*(9-6);"
assert 4 "(3+5)/2 ;"
assert 10 "-10+20;"
assert 5 "+2+4-1;"

assert 1 "1 == 1;"
assert 0 "1 != 1;"

assert 0 "2 == 1;"
assert 1 "2 != 1;"

assert 1 "1 < 2;"
assert 0 "2 < 1;"

assert 1 "1 <= 2;"
assert 1 "1 <= 1;"
assert 0 "2 <= 1;"

assert 0 "1 >= 2;"
assert 1 "1 >= 1;"
assert 1 "2 >= 1;"

assert 0 "(-1+2*20) >= (2*10/2+10+50);"
assert 1 "(5*5*2) == (10/2+50-5);"

assert 1 "a=1;a;"
assert 14 "a = 3; b=5 * 6-8;  a+b/2;"

assert 1 "aaaa=1;aaaa;"
assert 14 "abb = 3; aaaa=5 * 6-8;  abb+aaaa/2;"

assert 1 "aaaa=1; return aaaa;"
assert 14 "abb = 3; aaaa=5 * 6-8;  return abb+aaaa/2;"
assert 1 "return aaaa=1;aaaa;"
assert 22 "abb = 3; return aaaa=5 * 6-8;  abb+aaaa/2;"

assert 12 "a=1; if(a == 1) return 12;"
assert 0 "a=1; if(a != 1) return 12; return 0;"

assert 12 "a=1; if(a == 1) return 12; else return 0;"
assert 0 "a=1; if(a != 1) return 12; else return 0;"

assert 10 "a=0; while(a<10) a = a+1; return a;"
assert 10 "a=0; while(a<100) if(a == 10) return a; else a = a+1;"

assert 10 "for(a=0;a<10;a=a+1) b=10;return a;"
assert 10 "a=0;for(;a<10;a=a+1) b=10;return a;"
assert 10 "a=0;for(;a<10;) a=a+1;return a;"
assert 10 "for(a=0;a<10;) a=a+1;return a;"
assert 10 "for(a=0;;) if(a<10) a=a+1;else return a;"

assert 111 "{a=0;a=a+1;a=a+10;a=a+100;return a;}"

echo OK