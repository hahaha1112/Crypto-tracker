// 全局变量
let allCoins = [];
let allExchanges = [];
let marketChart = null;
let historyChart = null;
let coinDetailChart = null;
let currentCoin = 'BTC';

// DOM元素加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 初始化
    initNavigation();
    initCharts();
    loadCoinsAndExchanges();
    
    // 注册事件监听
    registerEventListeners();
    
    // 首次加载数据
    refreshPrices();
    setInterval(refreshPrices, 60000); // 每分钟刷新一次价格
});

// 初始化导航
function initNavigation() {
    // 页面导航处理
    document.querySelectorAll('#navbarNav .nav-link').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href').substring(1);
            showPage(targetId);
            
            // 更新导航栏高亮
            document.querySelectorAll('#navbarNav .nav-link').forEach(navLink => {
                navLink.classList.remove('active');
            });
            this.classList.add('active');
        });
    });
}

// 显示指定页面
function showPage(pageId) {
    document.querySelectorAll('.content-page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
    
    // 针对不同页面执行特定操作
    if (pageId === 'dashboard') {
        refreshPrices();
    } else if (pageId === 'alerts') {
        loadAlerts();
    } else if (pageId === 'history') {
        updateHistoryChart();
    }
}

// 初始化图表
function initCharts() {
    // 市场概览图表
    const marketCtx = document.getElementById('marketOverviewChart').getContext('2d');
    marketChart = new Chart(marketCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: []
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: '主要加密货币价格走势'
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: '时间'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: '价格 (USD)'
                    },
                    beginAtZero: false
                }
            }
        }
    });
    
    // 历史价格图表
    const historyCtx = document.getElementById('historyChart').getContext('2d');
    historyChart = new Chart(historyCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: '价格',
                data: [],
                borderColor: '#007bff',
                backgroundColor: 'rgba(0, 123, 255, 0.1)',
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false
                },
                title: {
                    display: true,
                    text: '价格历史'
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: '时间'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: '价格 (USD)'
                    },
                    beginAtZero: false
                }
            }
        }
    });
    
    // 币种详情图表
    const coinDetailCtx = document.getElementById('coinDetailChart').getContext('2d');
    coinDetailChart = new Chart(coinDetailCtx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: '价格',
                data: [],
                borderColor: '#007bff',
                backgroundColor: 'rgba(0, 123, 255, 0.1)',
                fill: true,
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false
                },
                title: {
                    display: true,
                    text: '24小时价格走势'
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: '时间'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: '价格 (USD)'
                    },
                    beginAtZero: false
                }
            }
        }
    });
}

// 加载币种和交易所数据
function loadCoinsAndExchanges() {
    // 获取支持的加密货币列表
    fetch('/api/v1/coins')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                allCoins = data.data;
                updateCoinSelectors();
            }
        })
        .catch(error => console.error('获取币种列表失败:', error));
    
    // 获取支持的交易所列表
    fetch('/api/v1/exchanges')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                allExchanges = data.data;
                updateExchangeSelectors();
            }
        })
        .catch(error => console.error('获取交易所列表失败:', error));
}

// 更新币种选择下拉框
function updateCoinSelectors() {
    const coinSelectors = [
        document.getElementById('coinSelector'),
        document.getElementById('alertCoin')
    ];
    
    coinSelectors.forEach(selector => {
        if (selector) {
            selector.innerHTML = '';
            allCoins.forEach(coin => {
                const option = document.createElement('option');
                option.value = coin.symbol;
                option.textContent = `${coin.symbol} - ${coin.name}`;
                selector.appendChild(option);
            });
        }
    });
}

// 更新交易所选择下拉框
function updateExchangeSelectors() {
    const exchangeSelector = document.getElementById('alertExchange');
    if (exchangeSelector) {
        exchangeSelector.innerHTML = '';
        allExchanges.forEach(exchange => {
            const option = document.createElement('option');
            option.value = exchange;
            option.textContent = exchange;
            exchangeSelector.appendChild(option);
        });
    }
}

// 注册事件监听器
function registerEventListeners() {
    // 时间范围选择器变更
    const timeRangeSelector = document.getElementById('timeRangeSelector');
    if (timeRangeSelector) {
        timeRangeSelector.addEventListener('change', updateMarketChart);
    }
    
    // 历史图表相关选择器
    const coinSelector = document.getElementById('coinSelector');
    const historyRangeSelector = document.getElementById('historyRangeSelector');
    if (coinSelector) {
        coinSelector.addEventListener('change', function() {
            currentCoin = this.value;
            updateHistoryChart();
        });
    }
    if (historyRangeSelector) {
        historyRangeSelector.addEventListener('change', updateHistoryChart);
    }
    
    // 警报相关按钮
    const addAlertBtn = document.getElementById('addAlertBtn');
    if (addAlertBtn) {
        addAlertBtn.addEventListener('click', showAlertModal);
    }
    
    const saveAlertBtn = document.getElementById('saveAlertBtn');
    if (saveAlertBtn) {
        saveAlertBtn.addEventListener('click', saveAlert);
    }
    
    // 币种详情相关
    const createAlertForCoin = document.getElementById('createAlertForCoin');
    if (createAlertForCoin) {
        createAlertForCoin.addEventListener('click', function() {
            // 关闭当前模态框
            const coinDetailModal = bootstrap.Modal.getInstance(document.getElementById('coinDetailModal'));
            coinDetailModal.hide();
            
            // 打开警报模态框，并设置当前币种
            const alertCoinSelect = document.getElementById('alertCoin');
            if (alertCoinSelect) {
                alertCoinSelect.value = currentCoin;
            }
            
            showAlertModal();
        });
    }
}

// 刷新价格数据
function refreshPrices() {
    fetch('/api/v1/prices')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updatePricesTable(data.data);
                updateMarketChart();
            }
        })
        .catch(error => console.error('获取价格数据失败:', error));
}

// 更新价格表格
function updatePricesTable(pricesData) {
    const tableBody = document.querySelector('#pricesTable tbody');
    if (!tableBody) return;
    
    tableBody.innerHTML = '';
    
    // 假设pricesData是一个嵌套对象，交易所->币种->价格数据
    Object.keys(pricesData).forEach(exchange => {
        Object.keys(pricesData[exchange]).forEach(coinSymbol => {
            const price = pricesData[exchange][coinSymbol];
            const coin = allCoins.find(c => c.symbol === coinSymbol) || { name: coinSymbol, symbol: coinSymbol };
            
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${coin.symbol}</td>
                <td>${coin.name || coin.symbol}</td>
                <td>$${formatNumber(price.price)}</td>
                <td class="change ${price.change24h >= 0 ? 'positive' : 'negative'}">
                    ${price.change24h >= 0 ? '+' : ''}${price.change24h.toFixed(2)}%
                </td>
                <td>$${formatLargeNumber(price.volume24h)}</td>
                <td>
                    <button class="btn btn-sm btn-primary btn-action view-coin" data-coin="${coin.symbol}">
                        查看
                    </button>
                    <button class="btn btn-sm btn-outline-primary btn-action create-alert" data-coin="${coin.symbol}">
                        警报
                    </button>
                </td>
            `;
            
            // 添加事件监听
            const viewBtn = tr.querySelector('.view-coin');
            viewBtn.addEventListener('click', function() {
                showCoinDetail(coin.symbol);
            });
            
            const alertBtn = tr.querySelector('.create-alert');
            alertBtn.addEventListener('click', function() {
                const alertCoinSelect = document.getElementById('alertCoin');
                if (alertCoinSelect) {
                    alertCoinSelect.value = coin.symbol;
                }
                showAlertModal();
            });
            
            tableBody.appendChild(tr);
        });
    });
}

// 更新市场概览图表
function updateMarketChart() {
    const timeRange = document.getElementById('timeRangeSelector')?.value || '24h';
    
    // 这里应该从API获取历史数据，但由于是模拟实现，我们生成一些随机数据
    const topCoins = ['BTC', 'ETH', 'BNB', 'XRP'];
    const timeLabels = generateTimeLabels(timeRange);
    
    const datasets = topCoins.map((coin, index) => {
        const basePrice = getBasePriceForCoin(coin);
        const data = generateRandomPriceData(basePrice, timeLabels.length);
        const colors = [
            { border: '#007bff', background: 'rgba(0, 123, 255, 0.1)' },
            { border: '#28a745', background: 'rgba(40, 167, 69, 0.1)' },
            { border: '#ffc107', background: 'rgba(255, 193, 7, 0.1)' },
            { border: '#dc3545', background: 'rgba(220, 53, 69, 0.1)' }
        ];
        
        return {
            label: coin,
            data: data,
            borderColor: colors[index].border,
            backgroundColor: colors[index].background,
            fill: false,
            tension: 0.1
        };
    });
    
    marketChart.data.labels = timeLabels;
    marketChart.data.datasets = datasets;
    marketChart.update();
}

// 更新历史价格图表
function updateHistoryChart() {
    const coin = document.getElementById('coinSelector')?.value || 'BTC';
    const timeRange = document.getElementById('historyRangeSelector')?.value || '24h';
    const exchange = '模拟交易所';
    
    // 这里应该从API获取数据，但暂时用模拟数据
    const timeLabels = generateTimeLabels(timeRange);
    const basePrice = getBasePriceForCoin(coin);
    const priceData = generateRandomPriceData(basePrice, timeLabels.length);
    
    historyChart.data.labels = timeLabels;
    historyChart.data.datasets[0].data = priceData;
    historyChart.options.plugins.title.text = `${coin} 价格历史 (${timeRange})`;
    historyChart.update();
}

// 显示币种详情
function showCoinDetail(symbol) {
    currentCoin = symbol;
    const coin = allCoins.find(c => c.symbol === symbol) || { name: symbol, symbol: symbol };
    
    // 更新模态框内容
    document.getElementById('coinName').textContent = `${coin.name} (${coin.symbol})`;
    document.getElementById('coinLogo').src = coin.logoUrl || '/static/images/default-coin.png';
    document.getElementById('coinLogo').alt = coin.symbol;
    document.getElementById('coinDescription').textContent = coin.description || '暂无描述';
    
    // 获取当前价格
    fetch(`/api/v1/prices/${symbol}?exchange=模拟交易所`)
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const price = data.data;
                document.getElementById('coinPrice').textContent = `$${formatNumber(price.price)}`;
                
                const changeElement = document.getElementById('coinChange');
                changeElement.textContent = `${price.change24h >= 0 ? '+' : ''}${price.change24h.toFixed(2)}%`;
                changeElement.className = `change ${price.change24h >= 0 ? 'positive' : 'negative'}`;
                
                // 更新图表
                updateCoinDetailChart(symbol);
            }
        })
        .catch(error => console.error('获取币种价格失败:', error));
    
    // 显示模态框
    const coinDetailModal = new bootstrap.Modal(document.getElementById('coinDetailModal'));
    coinDetailModal.show();
}

// 更新币种详情图表
function updateCoinDetailChart(symbol) {
    // 获取24小时历史数据
    const timeLabels = generateTimeLabels('24h');
    const basePrice = getBasePriceForCoin(symbol);
    const priceData = generateRandomPriceData(basePrice, timeLabels.length);
    
    coinDetailChart.data.labels = timeLabels;
    coinDetailChart.data.datasets[0].data = priceData;
    coinDetailChart.options.plugins.title.text = `${symbol} 24小时价格走势`;
    coinDetailChart.update();
}

// 加载警报列表
function loadAlerts() {
    fetch('/api/v1/alerts')
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                updateAlertsTable(data.data);
            }
        })
        .catch(error => console.error('获取警报列表失败:', error));
}

// 更新警报表格
function updateAlertsTable(alerts) {
    const tableBody = document.querySelector('#alertsTable tbody');
    if (!tableBody) return;
    
    tableBody.innerHTML = '';
    
    alerts.forEach(alert => {
        const tr = document.createElement('tr');
        
        // 格式化条件类型
        let condition;
        switch (alert.type) {
            case 'price_above': condition = '价格高于'; break;
            case 'price_below': condition = '价格低于'; break;
            case 'change_above': condition = '涨幅高于'; break;
            case 'change_below': condition = '跌幅低于'; break;
            default: condition = alert.type;
        }
        
        // 格式化状态
        let statusClass;
        switch (alert.status) {
            case 'active': statusClass = 'status-active'; break;
            case 'triggered': statusClass = 'status-triggered'; break;
            case 'disabled': statusClass = 'status-disabled'; break;
            default: statusClass = '';
        }
        
        tr.innerHTML = `
            <td>${alert.coin}</td>
            <td>${alert.exchange}</td>
            <td>${condition}</td>
            <td>${alert.type.startsWith('price') ? '$' : ''}${alert.threshold}${alert.type.startsWith('change') ? '%' : ''}</td>
            <td><span class="status-badge ${statusClass}">${alert.status}</span></td>
            <td>${formatDate(new Date(alert.createdAt))}</td>
            <td>
                <button class="btn btn-sm btn-primary btn-action edit-alert" data-id="${alert.id}">
                    编辑
                </button>
                <button class="btn btn-sm btn-danger btn-action delete-alert" data-id="${alert.id}">
                    删除
                </button>
            </td>
        `;
        
        // 添加事件监听
        tr.querySelector('.edit-alert').addEventListener('click', function() {
            editAlert(alert.id);
        });
        
        tr.querySelector('.delete-alert').addEventListener('click', function() {
            deleteAlert(alert.id);
        });
        
        tableBody.appendChild(tr);
    });
}

// 显示警报模态框
function showAlertModal(alertId = null) {
    // 清空表单
    document.getElementById('alertId').value = alertId || '';
    document.getElementById('alertThreshold').value = '';
    document.getElementById('alertMessage').value = '';
    
    // 设置标题
    document.querySelector('#alertModal .modal-title').textContent = alertId ? '编辑价格警报' : '新建价格警报';
    
    // 如果是编辑模式，加载警报数据
    if (alertId) {
        fetch(`/api/v1/alerts/${alertId}`)
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    const alert = data.data;
                    document.getElementById('alertCoin').value = alert.coin;
                    document.getElementById('alertExchange').value = alert.exchange;
                    document.getElementById('alertType').value = alert.type;
                    document.getElementById('alertThreshold').value = alert.threshold;
                    document.getElementById('alertMessage').value = alert.message || '';
                }
            })
            .catch(error => console.error('获取警报详情失败:', error));
    }
    
    // 显示模态框
    const alertModal = new bootstrap.Modal(document.getElementById('alertModal'));
    alertModal.show();
}

// 保存警报
function saveAlert() {
    const alertId = document.getElementById('alertId').value;
    const alertData = {
        coin: document.getElementById('alertCoin').value,
        exchange: document.getElementById('alertExchange').value,
        type: document.getElementById('alertType').value,
        threshold: parseFloat(document.getElementById('alertThreshold').value),
        message: document.getElementById('alertMessage').value
    };
    
    // 验证数据
    if (!alertData.coin || !alertData.exchange || isNaN(alertData.threshold)) {
        alert('请填写完整的警报信息');
        return;
    }
    
    // 确定API端点和方法
    const url = alertId ? `/api/v1/alerts/${alertId}` : '/api/v1/alerts';
    const method = alertId ? 'PUT' : 'POST';
    
    // 发送请求
    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(alertData)
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // 关闭模态框
                const alertModal = bootstrap.Modal.getInstance(document.getElementById('alertModal'));
                alertModal.hide();
                
                // 刷新警报列表
                loadAlerts();
            } else {
                alert(`保存警报失败: ${data.error}`);
            }
        })
        .catch(error => console.error('保存警报失败:', error));
}

// 编辑警报
function editAlert(alertId) {
    showAlertModal(alertId);
}

// 删除警报
function deleteAlert(alertId) {
    if (confirm('确定要删除此警报吗？')) {
        fetch(`/api/v1/alerts/${alertId}`, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    loadAlerts();
                } else {
                    alert(`删除警报失败: ${data.error}`);
                }
            })
            .catch(error => console.error('删除警报失败:', error));
    }
}

// 生成时间标签
function generateTimeLabels(timeRange) {
    const now = new Date();
    const labels = [];
    let numPoints;
    let intervalMinutes;
    
    switch (timeRange) {
        case '24h':
            numPoints = 24;
            intervalMinutes = 60;
            break;
        case '7d':
            numPoints = 7;
            intervalMinutes = 60 * 24;
            break;
        case '30d':
            numPoints = 30;
            intervalMinutes = 60 * 24;
            break;
        default:
            numPoints = 24;
            intervalMinutes = 60;
    }
    
    for (let i = numPoints - 1; i >= 0; i--) {
        const time = new Date(now.getTime() - i * intervalMinutes * 60 * 1000);
        
        if (timeRange === '24h') {
            labels.push(time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }));
        } else {
            labels.push(time.toLocaleDateString([], { month: 'short', day: 'numeric' }));
        }
    }
    
    return labels;
}

// 获取币种基础价格
function getBasePriceForCoin(symbol) {
    switch (symbol) {
        case 'BTC': return 50000;
        case 'ETH': return 3000;
        case 'BNB': return 500;
        case 'XRP': return 1;
        case 'ADA': return 1.5;
        case 'SOL': return 100;
        case 'DOGE': return 0.1;
        case 'DOT': return 20;
        default: return 10;
    }
}

// 生成随机价格数据
function generateRandomPriceData(basePrice, numPoints) {
    const volatility = basePrice * 0.05; // 5%波动率
    const data = [];
    let currentPrice = basePrice;
    
    for (let i = 0; i < numPoints; i++) {
        // 添加一些随机波动
        const change = (Math.random() - 0.5) * volatility;
        currentPrice += change;
        if (currentPrice < 0) currentPrice = 0.01; // 防止价格为负
        data.push(currentPrice);
    }
    
    return data;
}

// 格式化数字（添加千位分隔符）
function formatNumber(num) {
    return new Intl.NumberFormat('zh-CN', { 
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(num);
}

// 格式化大数字（用K, M, B等后缀）
function formatLargeNumber(num) {
    if (num >= 1e9) {
        return (num / 1e9).toFixed(2) + 'B';
    } else if (num >= 1e6) {
        return (num / 1e6).toFixed(2) + 'M';
    } else if (num >= 1e3) {
        return (num / 1e3).toFixed(2) + 'K';
    }
    return num.toFixed(2);
}

// 格式化日期
function formatDate(date) {
    return date.toLocaleString('zh-CN', { 
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}
