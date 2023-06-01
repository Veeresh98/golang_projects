/*

Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].
*/


//solution 
package slice_problems

import (
		"fmt"
)

func twoSum(nums []int, target int) []int {

    for i := 0; i < len(nums)-1;i++{
        for j :=i+1;j < len(nums);j++{
            if (nums[i] + nums[j] == target){
                return []int{i, j}
            }else{
                fmt.Println("no indeces to add")
            }
        }
    }
    return []int{}
}

/*
Roman numerals are represented by seven different symbols: I, V, X, L, C, D and M.

Symbol       Value
I             1
V             5
X             10
L             50
C             100
D             500
M             1000
*/

func romanToInt(s string) int {
    roman := map[byte]int{
        'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
    }
    leng := len(s)
    if leng == 0 {
        return 0
    }
    if leng == 1 {
        return roman[s[0]]
    }
    sum := roman[s[leng - 1]]
    for i := leng - 2; i >= 0; i-- {
        if roman[s[i]] < roman[s[i+1]] {
            sum -= roman[s[i]]
        } else {
            sum += roman[s[i]]
        }
    }
    return sum
}



