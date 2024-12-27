import React, { useState } from 'react';
import Modal from 'react-modal';
import "../CSS/Modal.css"

const RenameModal = ({ isOpen, onRequestClose, onRename, itemId, itemType, item}) => {
  const [inputUser, setInputUser] = useState("");

  const handleInputUser = (event) => {
    setInputUser(event.target.value)
  }

  const handleComfirm = async () => {
    if (inputUser != "") {
      const token = localStorage.getItem('token');

      if (!token) throw new Error('Token not found');

      const url = itemType === "directory" ? `http://localhost:8040/directories/update/${itemId}` : `http://localhost:8040/files/update/${itemId}`

      const body = {
        name: inputUser,
        size: item.details.size,
        data_create: item.details.data_create,
        directory_id: item.details.directory_id,
        user_id: item.details.user_id
      }

      console.log(body)
      console.log(url)

      const response = await fetch(url, {
          method: 'PUT',
          headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(body)
      });

      console.log(response.text)
      if (!response.ok) throw new Error('Failed to rename');

      onRename()
    }
  }

  React.useEffect(() => {
      if (isOpen) {
          console.log("Open")
      }
    }, [isOpen]);
  

  return (
    <Modal
      isOpen={isOpen}
      onRequestClose={onRequestClose}
      contentLabel="Rename Item"
      className="modal-content"
      overlayClassName="modal-overlay"
    >
      <div className="modal-header">Как вы хотите переименовать {itemType === "directory" ? "директорию" : "файл"}?</div>
      <input className="moda-input" type='text' onChange={handleInputUser}></input>
      <div className="modal-footer">
      <button
        className="modal-button"
        onClick={handleComfirm}
        // disabled={!targetDirectory && modalCurrentDirectory?.directory_id !== null} // Разрешить выбор корня
      >
        Переименовать
      </button>
        <button className="modal-button" onClick={onRequestClose}>Отмена</button>
      </div>
    </Modal>
  );
};


export default RenameModal;
