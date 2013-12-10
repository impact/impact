impact.pdf: impact.tex impact.bib
	pdflatex impact.tex
	bibtex impact
	pdflatex impact.tex
clean:
	rm *.aux *.bbl *.blg *-blx.bib *.log *.out *.run.xml
