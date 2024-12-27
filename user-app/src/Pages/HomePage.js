import React, { useState, useEffect } from 'react';
import MoveModal from '../Components/MoveModal';
import RenameModal from '../Components/RenameModal';
import '../CSS/Home.css'; // Подключение CSS для стилизации

const HomePage = () => {
  const [items, setItems] = useState([]);
  const [selectedItem, setSelectedItem] = useState(null);
  const [currentDirectory, setCurrentDirectory] = useState(null);
  const [isMoveModalOpen, setIsMoveModalOpen] = useState(false);
  const [isRenameModalOpen, setIsRenameModalOpen] = useState(false);

  const getUserIdFromToken = () => {
    const token = localStorage.getItem('token'); // Имя токена в localStorage
    if (!token) {
      console.error('Токен не найден');
      return null;
    }

    try {
      const payload = JSON.parse(atob(token.split('.')[1])); // Декодируем payload из токена
      return payload.user_id || null; // Извлекаем ID пользователя
    } catch (error) {
      console.error('Ошибка при декодировании токена:', error);
      return null;
    }
  };

  const fetchDirectory = async (directoryId) => {
    try {
      setItems([]); // Очистка массива перед загрузкой данных
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Токен не найден');
        return;
      }

      const response = await fetch(`http://localhost:8040/directories/${directoryId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Ошибка при загрузке данных директории');
      }

      const data = await response.json();

      const formattedItems = [
        ...(data.Directories || []).map((dir) => ({
          id: dir.id,
          name: dir.name,
          type: 'directory',
          size: `${dir.size} KB`,
          details: dir,
        })),
        ...(data.MetadataFiles || []).map((file) => ({
          id: file.id,
          name: file.name,
          type: 'file',
          size: `${file.size} KB`,
          details: file,
        })),
      ];

      setItems(formattedItems);
      setSelectedItem(null); // Снимаем выделение при переходе в новую директорию
      setCurrentDirectory(data);
    } catch (error) {
      console.error('Ошибка загрузки данных директории:', error);
    }
  };

  useEffect(() => {
    console.log("Items loaded:", items);
  }, [items]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = localStorage.getItem('token'); // Получаем токен из localStorage

        if (!token) {
          console.error('Токен не найден');
          return;
        }

        const userId = getUserIdFromToken(); // Получаем ID пользователя из токена
        if (!userId) {
          console.error('Не удалось извлечь ID пользователя');
          return;
        }

        const response = await fetch(`http://localhost:8040/directories/home/user/${userId}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`, // Добавляем токен в заголовок
          },
        }); // Используем ID пользователя в запросе
        if (!response.ok) {
          throw new Error('Ошибка при загрузке данных');
        }
        const data = await response.json();

        const formattedItems = [
          ...(data.Directories || []).map((dir) => ({
            id: dir.id,
            name: dir.name,
            type: 'directory',
            size: `${dir.size} KB`,
            details: dir,
          })),
          ...(data.MetadataFiles || []).map((file) => ({
            id: file.id,
            name: file.name,
            type: 'file',
            size: `${file.size} KB`,
            details: file,
          })),
        ];

        setItems(formattedItems);
        setCurrentDirectory(data); // Устанавливаем корневую директорию
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
      }
    };

    fetchData();
  }, []);

  const handleItemClick = (item) => {
    setSelectedItem(item);
  };

  const handleItemDoubleClick = (item) => {
    if (item.type === 'directory') {
      fetchDirectory(item.id);
    }
  };

  const handleGoToParentDirectory = () => {
    console.log(currentDirectory)
    if (currentDirectory && currentDirectory.directory_id) {
      fetchDirectory(currentDirectory.directory_id);
    } else {
      console.log('Вы находитесь в корневой директории.');
    }
  };

  const handleDelete = async () => {
    if (selectedItem) {
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Токен не найден');
        return;
      }
      
      const response = await fetch(`http://localhost:5091/api/File/${selectedItem.id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Ошибка при удалении файла');
      }

      const url = selectedItem.type === "directory" ? `http://localhost:8040/directories/delete/${selectedItem.id}` : `http://localhost:8040/files/delete/${selectedItem.id}`
      const responseDel = await fetch(url, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`, // Добавляем токен в заголовок
        },
      });

      if (!responseDel.ok) {
        throw new Error('Ошибка при удалении файла');
      }

      setSelectedItem(null);
      setItems([])
      fetchDirectory(currentDirectory.id);
      console.log(currentDirectory)
    }
  };

  const handleBackgroundClick = () => {
    setSelectedItem(null); // Снимаем выделение с элемента
    if (currentDirectory) {
        console.log(currentDirectory)
        setCurrentDirectory(currentDirectory); // Обновляем информацию о текущей директории
    }
  };

  const handleMove = (targetDirectoryId) => {
    if (selectedItem) {
        console.log(`Moving ${selectedItem.name} to directory ${targetDirectoryId}`);
        setItems([])
        fetchDirectory(currentDirectory.id);
        setIsMoveModalOpen(false);
    }
  };

  const handleRename = () => {
    if (selectedItem) {
      setItems([])
      fetchDirectory(currentDirectory.id);
      setIsRenameModalOpen(false);
    }
  }

 
  const handleUploadFile = async (event) => {
    try {
      const file = event.target.files[0];
      if (!file) return;

      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Токен не найден');
        return;
      }

      
      const formData = new FormData();
      formData.append('file', event.target.files[0]);
      console.log([...formData.entries()]);
      

      const response = await fetch(`http://localhost:5091/api/File`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Ошибка при загрузке файла');
      }

      const uploadedFile = await response.json();
      console.log(uploadedFile)
      console.log(currentDirectory)

      const body = {
        id: uploadedFile.id,
        size: uploadedFile.size,
        data_create: uploadedFile.creationDate,
        name: uploadedFile.name,
        directory_id: currentDirectory.id,
        user_id: getUserIdFromToken()
      }

      console.log(body)

      const metadataResponse = await fetch(`http://localhost:8040/files/create`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(body)
      });

      if (!metadataResponse.ok) {
        throw new Error('Ошибка при загрузке файла в приложение пользователя');
      }

      setItems([])
      fetchDirectory(currentDirectory.id);
    } catch (error) {
      console.error('Ошибка загрузки файла:', error);
    }
  };

  const handleCreateDirectory = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Токен не найден');
        return;
      }

      const body = {
        name: 'Новая папка',
        size: 0,
        directory_id: currentDirectory?.id,
        user_id: getUserIdFromToken()
      }

      console.log(body)

      const response = await fetch(`http://localhost:8040/directories/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      const textResponse = await response.text();
      console.log('Ответ сервера:', textResponse);

      if (!response.ok) {
        throw new Error('Ошибка при создании директории');
      }

      setItems([])
      fetchDirectory(currentDirectory.id);
    } catch (error) {
      console.error('Ошибка создания директории:', error);
    }
  };

  const handleDownloadFile = async () => {
    if (!selectedItem || selectedItem.type !== 'file') {
      console.error('Не выбран файл для скачивания');
      return;
    }
  
    try {
      window.location.href = `http://localhost:5091/api/file/${selectedItem.id}`;
    } catch (error) {
      console.error('Ошибка скачивания файла:', error);
    }
  };
  

  return (
    <div className="home-container">
      <div className="file-list"  onClick={handleBackgroundClick}>
        <h2>Файлы</h2>
        {currentDirectory && currentDirectory.directory_id && (
          <div className="back-arrow" onClick={handleGoToParentDirectory}>
            ← Назад
          </div>
        )}
        <div className="items-grid">
          {items.map((item) => (
            <div
              key={item.id}
              className={`item ${selectedItem?.id === item.id && selectedItem?.type === item.type ? 'selected' : ''}`}
              onClick={(e) => {
                e.stopPropagation(); // Предотвращаем клик на фон
                handleItemClick(item);
              }}
              onDoubleClick={() => handleItemDoubleClick(item)}
            >
              <div className={`icon ${item.type}`}></div>
              <p className="item-name">{item.name}</p>
            </div>
          ))}
        </div>
        <div className="action-buttons">
          <label htmlFor="file-upload" className="upload-label">
            Загрузить файл
          </label>
          <input
            id="file-upload"
            type="file"
            style={{ display: 'none' }}
            onChange={handleUploadFile}
          />
          <button onClick={handleCreateDirectory}>Создать директорию</button>
        </div>
      </div>

      <div className="file-details">
        <h2>Информация</h2>
        {selectedItem ? (
          <div>
            <p><strong>Имя:</strong> {selectedItem.details.name}</p>
            <p><strong>Размер:</strong> {selectedItem.details.size} KB</p>
            <p><strong>Путь:</strong> {selectedItem.details.path}</p>
            <p><strong>Дата создания:</strong> {new Date(selectedItem.details.data_create).toLocaleString()}</p>
            {selectedItem.details.data_change && (
              <p><strong>Дата изменения:</strong> {new Date(selectedItem.details.data_change).toLocaleString()}</p>
            )}
            <p><strong>Тип:</strong> {selectedItem.type === 'directory' ? 'Директория' : 'Файл'}</p>
            <button onClick={handleDelete}>Удалить</button>
            <button onClick={() => setIsMoveModalOpen(true)}>Переместить</button>
            <button onClick={() => setIsRenameModalOpen(true)}>Переименовать</button>
            {selectedItem.type === 'file' && (
              <button onClick={handleDownloadFile}>Скачать файл</button>
            )}
          </div>
        )  : currentDirectory ? (
            <div>
              <p><strong>Имя директории:</strong> {currentDirectory.name}</p>
              <p><strong>Путь:</strong> {currentDirectory.path}</p>
              <p><strong>Количество файлов:</strong> {currentDirectory.count_files}</p>
              <p><strong>Количество поддиректорий:</strong> {currentDirectory.count_directories}</p>
              <p><strong>Дата создания:</strong> {new Date(currentDirectory.data_create).toLocaleString()}</p>
            </div>
          ) : (
          <p>Выберите элемент, чтобы увидеть информацию.</p>
        )}
      </div>

      <MoveModal
        isOpen={isMoveModalOpen}
        onRequestClose={() => setIsMoveModalOpen(false)}
        currentDirectoryId={currentDirectory?.id || null}
        onMove={handleMove}
        itemToMoveId={selectedItem?.id || null} 
        itemToMoveType={selectedItem?.type || null} 
      />

      <RenameModal 
        isOpen={isRenameModalOpen}
        onRequestClose={() => setIsRenameModalOpen(false)}
        onRename={handleRename}
        itemId={selectedItem?.id || null}
        itemType={selectedItem?.type || null}
        item={selectedItem}
      />
    </div>
  );
};

export default HomePage;
