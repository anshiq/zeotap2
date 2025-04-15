document.addEventListener('DOMContentLoaded', function() {
    // UI Elements
    const sourceTypeSelect = document.getElementById('sourceType');
    const clickhouseConfig = document.getElementById('clickhouseConfig');
    const flatfileConfig = document.getElementById('flatfileConfig');
    const tableSection = document.querySelector('.table-section');
    const columnsSection = document.querySelector('.columns-section');
    const previewSection = document.querySelector('.preview-section');
    const ingestionSection = document.querySelector('.ingestion-section');
    const directionSelect = document.getElementById('direction');
    const targetFileSection = document.getElementById('targetFileSection');
    const targetTableSection = document.getElementById('targetTableSection');
    
    // Buttons
    const connectBtn = document.getElementById('connectBtn');
    const loadFileBtn = document.getElementById('loadFileBtn');
    const loadColumnsBtn = document.getElementById('loadColumnsBtn');
    const previewBtn = document.getElementById('previewBtn');
    const startIngestionBtn = document.getElementById('startIngestionBtn');
    
    // Status elements
    const statusMessage = document.getElementById('statusMessage');
    const resultMessage = document.getElementById('resultMessage');
    
    // Event listeners
    sourceTypeSelect.addEventListener('change', function() {
        if (this.value === 'clickhouse') {
            clickhouseConfig.style.display = 'block';
            flatfileConfig.style.display = 'none';
        } else {
            clickhouseConfig.style.display = 'none';
            flatfileConfig.style.display = 'block';
        }
        hideSections();
    });
    
    directionSelect.addEventListener('change', function() {
        if (this.value === 'clickhouse_to_flatfile') {
            targetFileSection.style.display = 'block';
            targetTableSection.style.display = 'none';
        } else {
            targetFileSection.style.display = 'none';
            targetTableSection.style.display = 'block';
        }
    });
    
    connectBtn.addEventListener('click', connectToSource);
    loadFileBtn.addEventListener('click', loadFile);
    loadColumnsBtn.addEventListener('click', loadColumns);
    previewBtn.addEventListener('click', previewData);
    startIngestionBtn.addEventListener('click', startIngestion);
    
    // File upload handler
    document.getElementById('fileUpload').addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            document.getElementById('ffFilePath').value = file.name;
            // In a real app, you would handle the file upload here
        }
    });
    
    function hideSections() {
        tableSection.style.display = 'none';
        columnsSection.style.display = 'none';
        previewSection.style.display = 'none';
        ingestionSection.style.display = 'none';
    }
    
    function showStatus(message, isError = false) {
        statusMessage.textContent = message;
        statusMessage.className = isError ? 'error' : 'info';
    }
    
    function showResult(message, isError = false) {
        resultMessage.textContent = message;
        resultMessage.className = isError ? 'error' : 'success';
    }
    
    function connectToSource() {
        const sourceType = sourceTypeSelect.value;
        showStatus(`Connecting to ${sourceType}...`);
        
        let config = {};
        if (sourceType === 'clickhouse') {
            config = {
                host: document.getElementById('chHost').value,
                port: document.getElementById('chPort').value,
                database: document.getElementById('chDatabase').value,
                user: document.getElementById('chUser').value,
                jwtToken: document.getElementById('chJWTToken').value,
                secure: document.getElementById('chSecure').checked
            };
        } else {
            config = {
                filePath: document.getElementById('ffFilePath').value,
                delimiter: document.getElementById('ffDelimiter').value || ','
            };
        }
        
        fetch('/api/connect', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                sourceType: sourceType,
                clickHouseConfig: sourceType === 'clickhouse' ? config : {},
                flatFileConfig: sourceType === 'flatfile' ? config : {}
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                throw new Error(data.error);
            }
            showStatus('Connection successful!');
            showResult('');
            
            // Show tables section
            tableSection.style.display = 'block';
            
            if (sourceType === 'clickhouse') {
                // Load tables for ClickHouse
                loadClickHouseTables();
            } else {
                // For flat file, we just have one "table"
                const tableSelect = document.getElementById('tableSelect');
                tableSelect.innerHTML = '<option value="file_data">File Data</option>';
            }
        })
        .catch(error => {
            showStatus('');
            showResult(`Connection failed: ${error.message}`, true);
        });
    }
    
    function loadClickHouseTables() {
        const sourceType = sourceTypeSelect.value;
        showStatus('Loading tables...');
        
        const params = new URLSearchParams();
        params.append('sourceType', sourceType);
        params.append('host', document.getElementById('chHost').value);
        params.append('port', document.getElementById('chPort').value);
        params.append('database', document.getElementById('chDatabase').value);
        params.append('user', document.getElementById('chUser').value);
        params.append('jwtToken', document.getElementById('chJWTToken').value);
        params.append('secure', document.getElementById('chSecure').checked);
        
        fetch(`/api/tables?${params.toString()}`)
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                throw new Error(data.error);
            }
            
            const tableSelect = document.getElementById('tableSelect');
            tableSelect.innerHTML = '';
            
            data.tables.forEach(table => {
                const option = document.createElement('option');
                option.value = table;
                option.textContent = table;
                tableSelect.appendChild(option);
            });
            
            showStatus('');
        })
        .catch(error => {
            showStatus('');
            showResult(`Failed to load tables: ${error.message}`, true);
        });
    }
    
    function loadFile() {
        // In a real implementation, this would handle file upload
        // For this example, we'll just proceed to show the table section
        tableSection.style.display = 'block';
        showStatus('File loaded successfully!');
    }
    
    function loadColumns() {
        const sourceType = sourceTypeSelect.value;
        const table = document.getElementById('tableSelect').value;
        showStatus('Loading columns...');
        
        const params = new URLSearchParams();
        params.append('sourceType', sourceType);
        params.append('table', table);
        
        if (sourceType === 'clickhouse') {
            params.append('host', document.getElementById('chHost').value);
            params.append('port', document.getElementById('chPort').value);
            params.append('database', document.getElementById('chDatabase').value);
            params.append('user', document.getElementById('chUser').value);
            params.append('jwtToken', document.getElementById('chJWTToken').value);
            params.append('secure', document.getElementById('chSecure').checked);
        } else {
            params.append('filePath', document.getElementById('ffFilePath').value);
            params.append('delimiter', document.getElementById('ffDelimiter').value || ',');
        }
        
        fetch(`/api/columns?${params.toString()}`)
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                throw new Error(data.error);
            }
            
            const columnsList = document.getElementById('columnsList');
            columnsList.innerHTML = '';
            
            data.columns.forEach(column => {
                const div = document.createElement('div');
                div.className = 'column-item';
                
                const checkbox = document.createElement('input');
                checkbox.type = 'checkbox';
                checkbox.id = `col_${column.name}`;
                checkbox.value = column.name;
                checkbox.checked = true;
                
                const label = document.createElement('label');
                label.htmlFor = `col_${column.name}`;
                label.textContent = `${column.name} (${column.type})`;
                
                div.appendChild(checkbox);
                div.appendChild(label);
                columnsList.appendChild(div);
            });
            
            columnsSection.style.display = 'block';
            ingestionSection.style.display = 'block';
            showStatus('');
        })
        .catch(error => {
            showStatus('');
            showResult(`Failed to load columns: ${error.message}`, true);
        });
    }
    
    function previewData() {
        const sourceType = sourceTypeSelect.value;
        const table = document.getElementById('tableSelect').value;
        showStatus('Loading preview...');
        
        // Get selected columns
        const selectedColumns = [];
        document.querySelectorAll('#columnsList input[type="checkbox"]:checked').forEach(checkbox => {
            selectedColumns.push(checkbox.value);
        });
        
        let config = {};
        if (sourceType === 'clickhouse') {
            config = {
                host: document.getElementById('chHost').value,
                port: document.getElementById('chPort').value,
                database: document.getElementById('chDatabase').value,
                user: document.getElementById('chUser').value,
                jwtToken: document.getElementById('chJWTToken').value,
                secure: document.getElementById('chSecure').checked
            };
        } else {
            config = {
                filePath: document.getElementById('ffFilePath').value,
                delimiter: document.getElementById('ffDelimiter').value || ','
            };
        }
        
        fetch('/api/preview', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                sourceType: sourceType,
                table: table,
                columns: selectedColumns,
                clickHouseConfig: sourceType === 'clickhouse' ? config : {},
                flatFileConfig: sourceType === 'flatfile' ? config : {}
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                throw new Error(data.error);
            }
            
            const previewData = document.getElementById('previewData');
            previewData.innerHTML = '';
            
            if (data.data.length === 0) {
                previewData.textContent = 'No data found';
                return;
            }
            
            // Create table
            const table = document.createElement('table');
            
            // Create header
            const thead = document.createElement('thead');
            const headerRow = document.createElement('tr');
            
            Object.keys(data.data[0]).forEach(key => {
                const th = document.createElement('th');
                th.textContent = key;
                headerRow.appendChild(th);
            });
            
            thead.appendChild(headerRow);
            table.appendChild(thead);
            
            // Create body
            const tbody = document.createElement('tbody');
            
            data.data.forEach(row => {
                const tr = document.createElement('tr');
                
                Object.values(row).forEach(value => {
                    const td = document.createElement('td');
                    td.textContent = value !== null ? value.toString() : 'NULL';
                    tr.appendChild(td);
                });
                
                tbody.appendChild(tr);
            });
            
            table.appendChild(tbody);
            previewData.appendChild(table);
            
            previewSection.style.display = 'block';
            showStatus('');
        })
        .catch(error => {
            showStatus('');
            showResult(`Failed to load preview: ${error.message}`, true);
        });
    }
    
    function startIngestion() {
        const sourceType = sourceTypeSelect.value;
        const table = document.getElementById('tableSelect').value;
        const direction = directionSelect.value;
        showStatus('Starting ingestion...');
        
        // Get selected columns
        const selectedColumns = [];
        document.querySelectorAll('#columnsList input[type="checkbox"]:checked').forEach(checkbox => {
            selectedColumns.push(checkbox.value);
        });
        
        let config = {};
        if (sourceType === 'clickhouse') {
            config = {
                host: document.getElementById('chHost').value,
                port: document.getElementById('chPort').value,
                database: document.getElementById('chDatabase').value,
                user: document.getElementById('chUser').value,
                jwtToken: document.getElementById('chJWTToken').value,
                secure: document.getElementById('chSecure').checked
            };
        } else {
            config = {
                filePath: document.getElementById('ffFilePath').value,
                delimiter: document.getElementById('ffDelimiter').value || ','
            };
        }
        
        const targetConfig = {
            filePath: document.getElementById('targetFilePath').value,
            delimiter: document.getElementById('ffDelimiter').value || ','
        };
        
        if (direction === 'flatfile_to_clickhouse') {
            targetConfig.database = document.getElementById('chDatabase').value;
            targetConfig.table = document.getElementById('targetTable').value;
        }
        
        fetch('/api/ingest', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                direction: direction,
                table: table,
                columns: selectedColumns,
                clickHouseConfig: config,
                flatFileConfig: targetConfig
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                throw new Error(data.error);
            }
            
            showStatus('Ingestion completed!');
            showResult(`Successfully processed ${data.count} records.`);
        })
        .catch(error => {
            showStatus('');
            showResult(`Ingestion failed: ${error.message}`, true);
        });
    }
});