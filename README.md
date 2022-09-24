# Eostre

### Status

- Working on reporting Task errors to the worker
- Working on a simple SQS client
- Working on tests :)

### Motivation

I hate lock-in!  

Finding a good library, writing thousands of lines, then having to drop the library because it doesn't support a very specific use case is the worst case scenario. Especially for something like task queueing. 

On the other hand writing your own task queue is a great pain.

Eostre is a task queueing library that hopes to be completely modular. From the worker, task manager, to the annoying reflect stuff, take what you want and nothing more.
