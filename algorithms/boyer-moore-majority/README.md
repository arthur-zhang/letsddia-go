# Boyer–Moore majority vote algorithm

Boyer-Moore majority vote algorithm is a technique used to find the majority element in an array efficiently with a time
complexity of O(n), where n is the length of the array. The majority element is defined as the element that appears more
than ⌊n/2⌋ times in the array.

Suppose we have an array [2, 4, 7, 4, 2, 4, 4]. Our goal is to find the majority element in this array.

1. Initialization: We start by assuming the first element of the array as the majority element and assign a count
   variable as 1.
2. Iteration: We iterate through the remaining elements of the array. For each element, we perform the following steps:

    - a. If the count variable is 0, we assign the current element as the assumed majority element and set the count to
        1.
    - b. If the current element is the same as the assumed majority element, we increment the count.
    - c. If the current element is different from the assumed majority element, we decrement the count.
3. Updating majority element: As we iterate through the array, if the count reaches 0, we update the assumed majority
   element to the current element and reset the count to 1.
4. Final check: After iterating through the entire array, the assumed majority element will be the potential majority
   element. We need to perform a final verification by counting its occurrences in the array.

In our example:

1. Initialize: Assume the majority element as 2 and set count as 1.

2. Iteration:
    - 4 is different from the assumed majority element (2), so decrement the count to 0.
    - Update the assumed majority element to 4 and reset the count to 1.
    - 7 is different from the assumed majority element (4), so decrement the count to 0.
    - Update the assumed majority element to 7 and reset the count to 1.
    - 4 is different from the assumed majority element (7), so decrement the count to 0.
    - Update the assumed majority element to 4 and reset the count to 1.
    - 4 is the same as the assumed majority element (4), so increment the count to 2.
    - 4 is the same as the assumed majority element (4), so increment the count to 3.
3. Final check: The assumed majority element after the iteration is 4. We need to verify if it appears more than
   ⌊n/2⌋ times in the array. Counting the occurrences, we find that 4 appears 4 times, which is more than ⌊7/2⌋ = 3.
   Therefore, the majority element in the array is 4.

The Boyer-Moore majority vote algorithm efficiently finds the majority element in linear time by canceling out
non-majority elements. It is particularly useful when there is guaranteed to be a majority element present in the
array.