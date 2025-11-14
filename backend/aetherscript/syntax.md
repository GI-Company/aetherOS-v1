# AetherScript Syntax

AetherScript uses a simple, LISP-like syntax. The basic building block of an AetherScript program is the S-expression (symbolic expression).

## S-Expressions

An S-expression is either an atom or a list of S-expressions. An atom is a number, a string, or a symbol. A list is a sequence of S-expressions enclosed in parentheses.

For example:

```
(print "Hello, world!")
(let x 10)
(print x)
```

## Atoms

### Numbers

Numbers can be integers or floating-point numbers.

### Strings

Strings are enclosed in double quotes.

### Symbols

Symbols are used to represent variables and function names.

## Lists

Lists are used to represent function calls and other language constructs.

## Comments

Comments start with a semicolon (;) and continue to the end of the line.

For example:

```
; This is a comment
(print "Hello, world!") ; This is also a comment
```
