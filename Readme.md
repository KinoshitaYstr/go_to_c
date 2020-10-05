### go de c compiler tsukutte miru zo !!
* https://www.sigbus.info/compilerbook

## BNF
* program = stmt*
* stmt = "{" stmt* "}" | expr ";" | "return" expr ";" | "if" "(" expr ")" stmt ( "else" stmt )? | "while" "(" expr ")" stmt | "for" "(" expr? ";" expr? ";" expr? ")" stmt
* expr = assign
* assign = equality ( "=" assign )?
* equality = relational ( "==" relational | "!=" relational )*
* relational = add ( "<" add | "<=" add | ">" add | ">=" add )*
* add = mul ( "+" mul | "-" mul )*
* mul = unary ( "\*" unary | "/" unary )*
* unary = ( "+" | "-" )? primary
* primary = num | ident | "(" expr ")"