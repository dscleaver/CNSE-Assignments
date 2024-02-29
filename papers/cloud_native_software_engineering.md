# Cloud Native Software Engineering Review

## David Cleaver

I found the section on polyglot programming to be interesting. I've been a
polyglot programmer for the majority of my career. While I spent a large portion
of that career primarily working in Java. I've also worked in Scheme, Prolog,
Scala, Go, and Python. I've always felt that learning and developing a habit of
learning other languages builds stronger programming skills across all
languages. I've met engineers in my career that didn't believe that they needed
to learn any language other than Java.

In the past I've based my arguments for using new programming languages based
largely on the specific concepts and style of programming supported in the
language. This makes it a hard sell for many of the reasons stated in the paper.
I've only been successful in moving one team in Comcast from Java to Scala, and
a new team inheriting the project wanted to take the time to rewrite it back in
Java.

The paper opened my eyes to a new dimension of the programming language argument
that I had not previously considered. Container size is a huge problem that
we've battled with the adoption of Concourse within the company and will only
continue to be relevant as we deploy more container workloads. Additionally, I
had never considered choosing a language that better supports fast startup times
or even basing selection on the SDK support for the cloud provider.

It's important that we as a company create an engineering community that can
better adapt and adopt languages based on cloud native considerations. These
engineers will be able to more agilely support workloads that they create and
inherit.
