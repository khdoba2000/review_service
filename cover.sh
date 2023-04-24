#!/bin/bash

# # Get the coverage percentage
coverage=$(go test -coverprofile=coverage.out ./... | grep -o '[0-9]*.[0-9]*%' | sed 's/%//')    
rm coverage.out


coverageSum=0
count=0

for c in $(echo "$coverage" | bc -l);
do 
    # echo $c
    coverageSum=$(echo $coverageSum+$c | bc -l)
    count=$(echo $count+1 | bc -l)
done

# echo $coverageSum
# echo $count

avgCoverage=$(echo $(echo $coverageSum/$count | bc -l))
echo $avgCoverage

if (( $(echo "$avgCoverage < 50" | bc -l) )); then
    echo "Coverage is less than 50%: $avgCoverage"
    exit 1
else 
    echo "Coverage is greater than or equal to 50%: $avgCoverage"
fi


#  # # Get the coverage percentage
#     coverage=$(go test -coverprofile=coverage.out ./... | grep -o '[0-9]*.[0-9]*%' | sed 's/%//')

#     # Check if coverage is less than 50%
#     if (( $(echo "$coverage < 50" | bc -l) )); then
#         echo "Coverage is less than 50%: $coverage"
#     else 
#         echo "Coverage is greater than or equal to 50%: $coverage"
#     fi

# rm coverage.out
