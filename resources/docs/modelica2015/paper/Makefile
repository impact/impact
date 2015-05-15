all: clean impact.pdf

impact.pdf: impact.tex impact.bib
	pdflatex impact.tex
	biber impact
	pdflatex impact.tex

clean:
	rm -f *.aux *.bbl *.blg *-blx.bib *.log *.out *.run.xml *.toc *.bcf *~
	rm -f *.nlo *.ist *.nls *.ilg

distclean: clean
	rm -f impact.pdf
	rm -rf auto/ _region_*
