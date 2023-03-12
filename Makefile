# Instalar las dependencias del requirements.txt

install:
	pip install -r requirements.txt
	
install_conda_req:
	conda install -c conda-forge gdal