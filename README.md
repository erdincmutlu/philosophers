
### Dining Philosophers Problem

My solution to dining philosophers problem. Algorithm is:

1. Pick the left fork. No timeout here, wait until manage to pick the left fork
2. Try to pick right fork with 3 seconds timeout.
3. If cannot pick the right fork, release the left fork, go back to step 1
4. If can pick the right fork, eat for random duration, upto 6 seconds.
5. Release both forks
6. Go back to step 1