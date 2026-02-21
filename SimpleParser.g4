/*
 * Parser для упрощенного подмножества Go
 * Поддерживает: переменные, базовые типы, операции, if-else, for, fmt.Println
 */

parser grammar SimpleParser;

options {
    tokenVocab = SimpleLexer;
}

// Program entry point
program
    : (declaration | statement)* EOF
    ;

// Variable declarations
declaration
    : VAR IDENTIFIER type_ (ASSIGN expression)? SEMICOLON
    | VAR IDENTIFIER ASSIGN expression SEMICOLON
    ;

type_
    : INT
    | STRING
    | BOOL
    ;

// Statements
statement
    : assignment SEMICOLON
    | ifStatement
    | forStatement
    | functionCall SEMICOLON
    | block
    ;

assignment
    : IDENTIFIER ASSIGN expression
    | IDENTIFIER DECLARE expression
    ;

// If-else statement
ifStatement
    : IF expression block (ELSE block)?
    ;

// For loop (simple syntax)
forStatement
    : FOR (assignment | condition)? SEMICOLON (condition)? SEMICOLON (assignment)? block
    | FOR condition block
    ;

condition
    : expression
    ;

// Block
block
    : LBRACE (statement)* RBRACE
    ;

// Function call (only fmt.Println)
functionCall
    : PRINTLN LPAREN argList RPAREN
    | PRINTLN LPAREN RPAREN
    ;

// Expression
expression
    : primary
    | NOT expression
    | MINUS expression
    | expression (STAR | DIV) expression
    | expression (PLUS | MINUS) expression
    | expression (EQ | NEQ | LT | GT | LE | GE) expression
    | expression AND expression
    | expression OR expression
    ;

primary
    : INT_LIT
    | STRING_LIT
    | IDENTIFIER
    | TRUE
    | FALSE
    | LPAREN expression RPAREN
    ;

// Argument list for function calls
argList
    : expression (COMMA expression)*
    ;
