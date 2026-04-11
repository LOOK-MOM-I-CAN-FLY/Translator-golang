/*
 * Parser для упрощенного подмножества Go
 * Поддерживает: переменные, базовые типы, операции, if-else, for, fmt.Println, структуры
 */

parser grammar SimpleParser;

options {
    tokenVocab = SimpleLexer;
}

// Program entry point
program
    : (typeDeclaration | declaration | statement)* EOF
    ;

// Type declarations: structs
typeDeclaration
    : structDeclaration
    ;

structDeclaration
    : TYPE IDENTIFIER STRUCT LBRACE structField* RBRACE SEMICOLON?
    | STRUCT IDENTIFIER LBRACE structField* RBRACE SEMICOLON?
    ;

structField
    : IDENTIFIER type_ SEMICOLON
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
    | IDENTIFIER
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
    : lvalue ASSIGN expression
    | IDENTIFIER DECLARE expression
    ;

lvalue
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

// If-else statement
ifStatement
    : IF expression block (ELSE block)?
    ;

// For loop (simple syntax)
forStatement
    : FOR forClause block
    | FOR expression block
    | FOR block
    ;

forClause
    : init=assignment? SEMICOLON cond=expression? SEMICOLON post=assignment?
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

// Expression with precedence
expression
    : orExpr
    ;

orExpr
    : andExpr (OR andExpr)*
    ;

andExpr
    : compExpr (AND compExpr)*
    ;

compExpr
    : addExpr (compOp addExpr)*
    ;

addExpr
    : mulExpr (addOp mulExpr)*
    ;

mulExpr
    : unaryExpr (mulOp unaryExpr)*
    ;

compOp
    : EQ | NEQ | LT | GT | LE | GE
    ;

addOp
    : PLUS | MINUS
    ;

mulOp
    : STAR | DIV
    ;

unaryExpr
    : (NOT | MINUS) unaryExpr
    | primary
    ;

primary
    : literal
    | IDENTIFIER (DOT IDENTIFIER)*
    | LPAREN expression RPAREN
    ;

literal
    : INT_LIT
    | STRING_LIT
    | TRUE
    | FALSE
    ;

// Argument list for function calls
argList
    : expression (COMMA expression)*
    ;
