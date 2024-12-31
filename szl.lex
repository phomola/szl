#############################################
### Silesian morphology - for experiments ###
#############################################

### paradigms

## substantives

# m1 chłop
# m2 pies
# m3 dōmb
# m4 gorol
# m5 kóń
# m6 nōż
# m7 gazda

# f1 ryba
# f2 studnia
# f3 kuchyń
# f4 kość

# n1 miasto
# n2 pole
# n3 kozani
# n4 cielã

## adjectives

# a1 dobry
# a2 głupi

# stems

@chłop nm1 chłop
@pies nm2 p{ie}s
@lōd nm3 l{ō}d
@hołda nf1 hołd
@wysoki a2 wysok

# paradigms

-nm1 nmh1s 0
-nm1 nmh2s a
-nm1 nmh3s u
-nm1 nmh4s a
-nm1 nmh5s 'e
-nm1 nmh6s 'e
-nm1 nmh7s ym
-nm1 nmh1p 'i
-nm1 nmh2p ōw
-nm1 nmh3p ōm
-nm1 nmh4p ōw
-nm1 nmh5p 'i
-nm1 nmh6p ach
-nm1 nmh7p ami

-nm2 nm1s 0
-nm2 nm2s a >ō,o >ie,0
-nm2 nm3s u >ō,o >ie,0
-nm2 nm4s a >ō,o >ie,0
-nm2 nm5s 'e >ō,o >ie,0
-nm2 nm6s 'e >ō,o >ie,0
-nm2 nm7s ym >ō,o >ie,0
-nm2 nm1p y >ō,o >ie,0
-nm2 nm2p ōw >ō,o >ie,0
-nm2 nm3p ōm >ō,o >ie,0
-nm2 nm4p y >ō,o >ie,0
-nm2 nm5p y >ō,o >ie,0
-nm2 nm6p ach >ō,o >ie,0
-nm2 nm7p ami >ō,o >ie,0

-nm3 nm1s 0
-nm3 nm2s a >ō,o
-nm3 nm3s u >ō,o
-nm3 nm4s 0
-nm3 nm5s 'e >ō,o
-nm3 nm6s 'e >ō,o
-nm3 nm7s ym >ō,o
-nm3 nm1p y >ō,o
-nm3 nm2p ōw >ō,o
-nm3 nm3p ōm >ō,o
-nm3 nm4p y >ō,o
-nm3 nm5p y >ō,o
-nm3 nm6p ach >ō,o
-nm3 nm7p ami >ō,o

-nf1 nf1s a
-nf1 nf2s y
-nf1 nf3s 'e
-nf1 nf4s ã
-nf1 nf5s o
-nf1 nf6s 'e
-nf1 nf7s ōm
-nf1 nf1p y
-nf1 nf2p 0
-nf1 nf3p ōm
-nf1 nf4p y
-nf1 nf5p y
-nf1 nf6p ach
-nf1 nf7p ami

-a2 af1s 'ŏ
-a2 af2s ij
-a2 af3s ij
-a2 af4s 'õ
-a2 af5s 'ŏ
-a2 af6s ij
-a2 af7s 'ōm
-a2 af1p i
-a2 af2p ich
-a2 af3p im
-a2 af4p i
-a2 af5p i
-a2 af6p ich
-a2 af7p imi

# replacements

!> p'e pie
!> d'e dzie
!> s'e sie
!> p'i pi
!> k'ŏ kŏ
!> k'õ kõ
!> k'ō kō

# feature structures

*nmh1s N [gender:"m",anim:"hum",case:"nom",number:"sg"] autosem
*nmh2s N [gender:"m",anim:"hum",case:"gen",number:"sg"] autosem
*nmh3s N [gender:"m",anim:"hum",case:"dat",number:"sg"] autosem
*nmh4s N [gender:"m",anim:"hum",case:"acc",number:"sg"] autosem
*nmh5s N [gender:"m",anim:"hum",case:"voc",number:"sg"] autosem
*nmh6s N [gender:"m",anim:"hum",case:"loc",number:"sg"] autosem
*nmh7s N [gender:"m",anim:"hum",case:"ins",number:"sg"] autosem
*nmh1p N [gender:"m",anim:"hum",case:"nom",number:"pl"] autosem
*nmh2p N [gender:"m",anim:"hum",case:"gen",number:"pl"] autosem
*nmh3p N [gender:"m",anim:"hum",case:"dat",number:"pl"] autosem
*nmh4p N [gender:"m",anim:"hum",case:"acc",number:"pl"] autosem
*nmh5p N [gender:"m",anim:"hum",case:"voc",number:"pl"] autosem
*nmh6p N [gender:"m",anim:"hum",case:"loc",number:"pl"] autosem
*nmh7p N [gender:"m",anim:"hum",case:"ins",number:"pl"] autosem

*nm1s N [gender:"m",anim:"nhum",case:"nom",number:"sg"] autosem
*nm2s N [gender:"m",anim:"nhum",case:"gen",number:"sg"] autosem
*nm3s N [gender:"m",anim:"nhum",case:"dat",number:"sg"] autosem
*nm4s N [gender:"m",anim:"nhum",case:"acc",number:"sg"] autosem
*nm5s N [gender:"m",anim:"nhum",case:"voc",number:"sg"] autosem
*nm6s N [gender:"m",anim:"nhum",case:"loc",number:"sg"] autosem
*nm7s N [gender:"m",anim:"nhum",case:"ins",number:"sg"] autosem
*nm1p N [gender:"m",anim:"nhum",case:"nom",number:"pl"] autosem
*nm2p N [gender:"m",anim:"nhum",case:"gen",number:"pl"] autosem
*nm3p N [gender:"m",anim:"nhum",case:"dat",number:"pl"] autosem
*nm4p N [gender:"m",anim:"nhum",case:"acc",number:"pl"] autosem
*nm5p N [gender:"m",anim:"nhum",case:"voc",number:"pl"] autosem
*nm6p N [gender:"m",anim:"nhum",case:"loc",number:"pl"] autosem
*nm7p N [gender:"m",anim:"nhum",case:"ins",number:"pl"] autosem

*nf1s N [gender:"f",case:"nom",number:"sg"] autosem
*nf2s N [gender:"f",case:"gen",number:"sg"] autosem
*nf3s N [gender:"f",case:"dat",number:"sg"] autosem
*nf4s N [gender:"f",case:"acc",number:"sg"] autosem
*nf5s N [gender:"f",case:"voc",number:"sg"] autosem
*nf6s N [gender:"f",case:"loc",number:"sg"] autosem
*nf7s N [gender:"f",case:"ins",number:"sg"] autosem
*nf1p N [gender:"f",case:"nom",number:"pl"] autosem
*nf2p N [gender:"f",case:"gen",number:"pl"] autosem
*nf3p N [gender:"f",case:"dat",number:"pl"] autosem
*nf4p N [gender:"f",case:"acc",number:"pl"] autosem
*nf5p N [gender:"f",case:"voc",number:"pl"] autosem
*nf6p N [gender:"f",case:"loc",number:"pl"] autosem
*nf7p N [gender:"f",case:"ins",number:"pl"] autosem

*af1s A [gender:"f",case:"nom",number:"sg"] autosem
*af2s A [gender:"f",case:"gen",number:"sg"] autosem
*af3s A [gender:"f",case:"dat",number:"sg"] autosem
*af4s A [gender:"f",case:"acc",number:"sg"] autosem
*af5s A [gender:"f",case:"voc",number:"sg"] autosem
*af6s A [gender:"f",case:"loc",number:"sg"] autosem
*af7s A [gender:"f",case:"ins",number:"sg"] autosem
*af1p A [gender:"f",case:"nom",number:"pl"] autosem
*af2p A [gender:"f",case:"gen",number:"pl"] autosem
*af3p A [gender:"f",case:"dat",number:"pl"] autosem
*af4p A [gender:"f",case:"acc",number:"pl"] autosem
*af5p A [gender:"f",case:"voc",number:"pl"] autosem
*af6p A [gender:"f",case:"loc",number:"pl"] autosem
*af7p A [gender:"f",case:"ins",number:"pl"] autosem
