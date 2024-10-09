#include <stdio.h>

int main() {
    int number;

    // Ask user for a number
    printf("Enter a number: ");
    scanf("%d", &number);

    // Branching based on the number
    if (number % 2 == 0) {
        // If number is even
        printf("The number %d is even.\n", number);
    } else {
        // If number is odd
        printf("The number %d is odd.\n", number);
    }

    return 0;
}
