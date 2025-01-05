# Harp

Harp is a statically typed, memory safe, procedural hobby programming language. Its interpreter is written in Go with minimal dependencies.

I'm working on this project after finishing the tree-walk interpreter section of [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom. As such, my interpreter is pretty similar. The main differences between Harp and Lox is that Harp is statically typed and doesn't support classes. I'd also like to really expand on the standard library and error reporting.

Once I finish the book, I'd like to come back and make a bytecode virtual machine like clox. If you haven't already checked this book out, I highly recommend it.

## The plan

I'll be writing the scanner and parser from scratch and then executing the syntax tree with a tree-walk interpreter.

This is my first time designing/implementing my own programming language, so a lot of the design/implementation decisions were made for the fun of it more than any real defensible reason. This language isn't mean to be particularly useful, but I tried to at least make sane decisions. The main inspiration is Go.

Features:

- Statically typed
- Garbage collected
- Supports structs but not classes

Steps:

- [ ] Scanner (in progress)
- [ ] Parser
- [ ] Interpreter
