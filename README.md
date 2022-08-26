# Dining Philosophers in Go

Two implementations of the Dining Philosophers problem using
two different  algorithms (well, three really, including a somewhat fake one :-). This is both a demo of the problem 
solution, and an example of the power and simplicity of Go's channels and go routines

## Running

Select the implementation using a command line arg:

    go run . <impl>

You can choose:
- `fingers` or `f` e.g
- `resourcehierarchy` or `rh`
- `chandymisra` or `cm`

or build it first:

    go build .; diningphils-go <impl>

## Algorithms

### Fingers

This is a toy solution where the Philosophers eat with their fingers - and therefore don't need forks!

### Resource Hierarchy

This is the classic resource hierarchy solution. Each resource (the forks) is given a priority (a partial order), and the algorithm
always acquires the resources in that order (loweset first). The forks are labeled 0 to N-1 around the table, with fork n being to the left of
philosopher n. Picking up the lowest fork first therefor means that philosophers 0 - N-2 will pick up their leftfork first, while philosopher
N - 1 (who has fork N - 1 to the left, and fork 0 to the right) will pick up the right fork first. This means
that the deadlock situation of all philosophers picking up their left (or right) forks first is avoided

This solution avoids deadlock, but is not fair. A philosopher who eats quicker than their neighbor will
get more than their fair share of spaghetti.

### Chandy-Misra algorithm

The Chandy-Misra solution relies on a distributed algorithm that allows philosophers to communicate with each other. See
(https://www.cs.utexas.edu/users/misra/scannedPdf.dir/DrinkingPhil.pdf). The idea here is to model the conflict on resources as a directed
graph, and prove that if such a graph is acyclic, no deadlocks will occur. The algorithm is proved correct
by proving that any state transformation that it produces does not make this graph acyclic.