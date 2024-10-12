document.getElementById('start-button').addEventListener('click', function() {
    const threadCount = document.getElementById('thread-count').value;
    const statusDiv = document.getElementById('status');
    statusDiv.innerHTML = 'Запрос данных...';

    fetch('/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ threadCount: parseInt(threadCount) })
    })
    .then(response => response.json())
    .then(data => {
        statusDiv.innerHTML = data.message;
    })
    .catch(error => {
        console.error('Ошибка:', error);
        statusDiv.innerHTML = 'Произошла ошибка.';
    });
});

document.getElementById('send-button').addEventListener('click', function() {
    const threadCount = document.getElementById('thread-count').value;
    const statusDiv = document.getElementById('status');
    statusDiv.innerHTML = 'Отправка данных...';

    fetch('/send', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ threadCount: parseInt(threadCount) })
    })
    .then(response => response.json())
    .then(data => {
        statusDiv.innerHTML = data.message;
    })
    .catch(error => {
        console.error('Ошибка:', error);
        statusDiv.innerHTML = 'Произошла ошибка.';
    });
});
