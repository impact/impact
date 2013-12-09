impact.pdf: impact.dvi
	dvipdf impact.dvi

impact.dvi: impact.tex
	latexmk impact.tex
