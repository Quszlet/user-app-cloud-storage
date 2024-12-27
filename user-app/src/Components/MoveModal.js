import React, { useState } from 'react';
import Modal from 'react-modal';
import "../CSS/Modal.css"

const MoveModal = ({ isOpen, onRequestClose, currentDirectoryId, onMove, itemToMoveId, itemToMoveType}) => {
  const [modalItems, setModalItems] = useState([]);
  const [modalCurrentDirectory, setModalCurrentDirectory] = useState(null);
  const [targetDirectory, setTargetDirectory] = useState(null);

  const fetchDirectoryForModal = async (directoryId) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Token not found');
  
      const response = await fetch(`http://localhost:8040/directories/${directoryId || ''}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });
  
      if (!response.ok) throw new Error('Failed to fetch directory data');
  
      const data = await response.json();
      const formattedItems = [
        ...(data.Directories || []).map((dir) => ({
          id: dir.id,
          name: dir.name,
          type: 'directory',
          details: dir,
        })).filter((dir) => !(itemToMoveType === 'directory' && dir.id === itemToMoveId)),
      ];
  
      setModalItems(formattedItems);
      setModalCurrentDirectory(data); // Обновляем текущую директорию
    } catch (error) {
      console.error('Error fetching modal directory:', error);
    }
  };

  const handleMoveConfirm = async () => {
    let targetDirId = targetDirectory ? targetDirectory.id : null; // Используем null для корня
    if (targetDirectory || modalCurrentDirectory?.directory_id === null) {
        const token = localStorage.getItem('token');
        if (!token) throw new Error('Token not found');

        // Доработать !!!
        if (itemToMoveType === "directory") {
          if (targetDirId == null) {
            targetDirId = modalCurrentDirectory.id
          }

          console.log(targetDirId)

          const body = {
            moved_obj_id: itemToMoveId,
            directory_id: targetDirId,
          }

          const responseMoveObj = await fetch(`http://localhost:8040/directories/move`, {
              method: 'PUT',
              headers: {
                  'Content-Type': 'application/json',
                  Authorization: `Bearer ${token}`,
              },
              body: JSON.stringify(body)
          });

          console.log(body)

          if (!responseMoveObj.ok) throw new Error('Failed to directory movee');

          const obj = await responseMoveObj.json() 
        } else {
            const body = {
                moved_obj_id: itemToMoveId,
                directory_id: targetDirId,
            }

            const responseMoveObj = await fetch(`http://localhost:8040/files/move`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(body)
            });

            if (!responseMoveObj.ok) throw new Error('Failed to fetch file get');
    
            const obj = await responseMoveObj.json() 
        
        }

    
        onMove(targetDirId);
    }
  };

  const handleItemDoubleClick = (item) => {
    if (item.type === 'directory') {
      fetchDirectoryForModal(item.id);
    }
  };

  React.useEffect(() => {
    if (isOpen) {
        console.log(currentDirectoryId)
        setTargetDirectory(null);
      fetchDirectoryForModal(currentDirectoryId);
    }
  }, [isOpen, currentDirectoryId]);

  return (
    <Modal
      isOpen={isOpen}
      onRequestClose={onRequestClose}
      contentLabel="Move Item"
      className="modal-content"
      overlayClassName="modal-overlay"
    >
      <div className="modal-header">Выберите директорию для перемещения</div>
      {modalCurrentDirectory && (
        <div>
          <button
            className="modal-button"
            onClick={() => fetchDirectoryForModal(modalCurrentDirectory?.directory_id)}
            disabled={!modalCurrentDirectory || modalCurrentDirectory.directory_id === null}
          >
            Наверх
          </button>
        </div>
      )}
      <ul className="modal-list">
        {modalItems.map((item) => (
          <li
            key={item.id}
            onClick={() => setTargetDirectory(item)}
            onDoubleClick={() => handleItemDoubleClick(item)}
            className={`modal-list-item ${targetDirectory?.id === item.id ? 'selected' : ''}`}
          >
            {item.name}
          </li>
        ))}
      </ul>
      <div className="modal-footer">
      <button
        className="modal-button"
        onClick={handleMoveConfirm}
        disabled={!targetDirectory && modalCurrentDirectory?.directory_id !== null} // Разрешить выбор корня
      >
        Переместить сюда
      </button>
        <button className="modal-button" onClick={onRequestClose}>Отмена</button>
      </div>
    </Modal>
  );
};


export default MoveModal;
