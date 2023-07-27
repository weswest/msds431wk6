# Overview of Project

This is an assignment for Week 6 in the Northwestern Masters in Data Science [MSDS-431 Data Engineering with Go](https://msdsgo.netlify.app/data-engineering-with-go/) course.

The purpose of this work is to test the impact of using concurrency on program execution.  In this exercise, we create the toy effort of producing OLS linear regressions on a given dataset, repeating the exercise 10,000 times.  In one incarnation, we execute this serially, one iteration at a time; in another incarnation, we execute concurrently.

### Key findings
1. The program is written to print out the regression results.  In testing, this was a material addition to program execution time, so I iterated on the assignment to produce results with and without output printing
2. In silent mode, the concurrent version was ~40% faster than the serial version.  So concurrency took a material chunk out of the total execution time
3. In silent mode, the concurrent version used ~3x more time from user processing and >10x more system processing.  Also not surprising, since concurrency carries material overhead for scheduling and thread processing
4. Verbose mode increased time to process across the board.  The general trends in 2&3 above were held
5. I guarantee my implementation of concurrency was sub-optimal.  I'm sure there's even more that could have been done to increase the performance difference.

### Recommendations
If we're using go, make sure we're keeping concurrency top of mind in our toolchest.  Most of the time we will value real time to completion over number of compute cycles, and concurrency allows us to better ake advantage of the relative costs and value.

### Other things to keep in mind

We were recommended to use the [gonum/stat](https://pkg.go.dev/gonum.org/v1/gonum/stat) package for running regressions, which has built-in functionality to execute simple linear regressions that produce results as y_hat = a + b*x.  We were to build regressions using a subset of the Boston Housing Study dataset (somewhat cleaned up to be less problematic), identifying predictors to forecast median home value.

The purpose of the exercise was to produce simple regressions many times, it was not to build a "best" model.  So while coefficients and R-Squared results can be produced, this exercise purposely did not adhere to statistical best practices (e.g., holdout samples for out of sample testing, covariance exploration, etc).

As always, this program was developed on a Mac although both Mac and Windows executables are provided.  This is because the Canvas website which manages assignments will only accept a .exe and won't accept a Mac executable.  The Mac executable has been tested and works; the Windows executable has not been tested.

# Program Structure and Use

The program is organized as follows:
```
.
+-- cmd
|   +-- concurrent
|       +-- concurrent.exe
|       +-- main.go
|       +-- main_test.go
|   +-- serial
|       +-- serial.exe
|       +-- main.go
|       +-- main_test.go
+-- data
|   +-- boston.csv
+-- go.mod
+-- go.sum
+-- README.md
```

Command line execution of the program assumes you are in the appropriate directory (concurrent or serial).  The iterations can be run silently by executing the program:


```bash
./concurrent
```

Or if you want to print out the results of the regressions, then call the verbose command:

```bash
./concurrent -verbose=verbose
```

# Results

### Execution Times

As noted in the key takeaways, the concurrent program finished sooner in real terms but at the cost of increased system usage.  All times in the table below are in seconds, and there were 10,000 iterations.

| Iteration                   | real    | user  | sys   | user+sys |
| --------------------------- | ------- | ----- | ----- | -------- |
| concurrent without printing | 0.039   | 0.171 | 0.057 | 0.228    |
| serial without printing     | 0.072   | 0.067 | 0.004 | 0.071    |
| Delta                       | \-0.033 | 0.104 | 0.053 | 0.157    |
|                             |         |       |       |          |
| concurrent with printing    | 0.198   | 0.387 | 0.133 | 0.52     |
| serial with printing        | 0.234   | 0.127 | 0.042 | 0.169    |
| Delta                       | \-0.036 | 0.26  | 0.091 | 0.351    |

### Regression Results

For what it's worth, the regression results used were as follows:

Crim vs Median Value: 24.03 + -0.42 * Crim , R-squared: 0.15
Rooms vs Median Value: -34.66 + 9.10 * Rooms , R-squared: 0.48

Crim is a measure of neighborhood criminality; Rooms is a measure of the number of rooms in the house for sale.

These results were cross-referenced with other statistical packages to ensure accuracy.

# FYI - assignment details motivating this work

### Management Problem

In data science, we speak of "autoML" and "grid search," thinking of systematic methods for finding best settings for hyperparameters in machine learning and neural network models. These kinds of problems, like stochastic local search in general, are well-suited for concurrent and parallel processing methods. In fact, whenever we need to train and test more than one model on a dataset, we can take advantage of concurrency. 

Methods that have traditionally described as computer-intensive, such as multi-fold cross-validation, bootstrap estimation, or any method that we affectionately describe as involving "brute force," are well-suited for concurrent and parallel processing. 

This assignment uses concurrent programming in Go to analyze data from a classic case study in applied statistics and machine learning: The Boston Housing Study.

### Assignment Requirements 

Take on the role of the independent contractor hired by the management consultancy to evaluate the potential advantages of concurrency in Go. In particular, the firm wants to see how much can be gained by using concurrency in training and testing machine learning models.

* Write a short Go program without concurrency that trains and tests two predictive models using the data from the Boston Housing Study. The models should predict the response variable mv (median value of homes in thousands of 1970 US dollars) from a subset of the explanatory variables, as described in Week 6 Assignment Data: Boston Housing Study.  
* Write another Go program that trains and tests the same models with concurrency using goroutines and channels.
* Run both programs 100 times, measuring CPU times. Which program runs faster?  Explain how you might expand on this mini experiment to show the advantages of concurrency in practice.
* In the README.md file of the repository, document fully the work completed, explaining your choice of modeling methods and comparing processing times with and without goroutines.  Advise management regarding potential gains associated with concurrency in Go.

Notes. Feel free to use any modeling package or library for the models as long as it has a Windows-compatible version. It is a good idea to choose packages that have numerous contributors and users, as well as recent activity.  We can find numerous packages under the Machine Learning section of Awesome Go Links to an external site.. The gonum Links to an external site. library also represents a possibility, as we can use what we know about matrix algebra to fit a least-squares linear regression. Or we could look into gonum/stat. With more than 100 contributors and more than 6 thousand stars, gonum is likely to stay around for a while. It is an active GitHub repository. Think of gonum as the Go alternative to Python numpy.

### Grading Guidelines (100 Total Points)

* Coding rules, organization, and aesthetics (20 points). Effective use of Go modules and idiomatic Go. Code should be readable, easy to understand. Variable and function names should be meaningful, specific rather than abstract. They should not be too long or too short. Avoid useless temporary variables and intermediate results. Code blocks and line breaks should be clean and consistent. Break large code blocks into smaller blocks that accomplish one task at a time. Utilize readable and easy-to-follow control flow (if/else blocks and for loops). Distribute the not rather than the switch (and/or) in complex Boolean expressions. Programs should be self-documenting, with comments explaining the logic behind the code (McConnell 2004, 777–817).
* Testing and software metrics (20 points). Employ unit tests of critical components, generating synthetic test data when appropriate. Generate program logs and profiles when appropriate. Monitor memory and processing requirements of code components and the entire program. If noted in the requirements definition, conduct a Monte Carlo performance benchmark.
* Design and development (20 points). Employ a clean, efficient, and easy-to-understand design that meets all aspects of the requirements definition and serves the use case. When possible, develop general-purpose code modules that can be reused in other programming projects.
* Documentation (20 points). Effective use of Git/GitHub, including a README.md Markdown file for each repository, noting the roles of programs and data and explaining how to test and use the application.
* Application (20 points). Delivery of an executable load module or application (.exe file for Windows or .app file for MacOS). The application should run to completion without issues. If user input is required, the application should check for valid/usable input and should provide appropriate explanation to users who provide incorrect input. The application should employ clean design for the user experience and user interface (UX/UI).

### Assignment Deliverables

* Text showing the link to the GitHub repository for the assignment
* README.md Markdown text file documentation for the assignment
* Zip archive of the GitHub repository
* Executable load module for the program/application (.exe for Windows or .app for MacOS)