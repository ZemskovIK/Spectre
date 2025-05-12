import React from "react";
import PropTypes from "prop-types";
import "./Modal.css";

const Modal = ({ title, text, onClose, backgroundImage }) => {
  return (
    <>
      <div className="modal-blur-overlay" />

      <div className="modal-container">
        <div
          className="modal-background"
          style={{ backgroundImage: `url(${backgroundImage})` }}
        />

        <div className="modal-content">
          <button className="close-button" onClick={onClose}>
            ×
          </button>

          <div className="modal-header">
            {" "}
            <h2 className="modal-title">Автор: {title}</h2>
          </div>

          <div className="modal-body">
            {" "}
            <p className="modal-text">{text}</p>
          </div>
        </div>
      </div>
    </>
  );
};

Modal.propTypes = {
  title: PropTypes.string.isRequired,
  text: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
  backgroundImage: PropTypes.string.isRequired,
};

export default Modal;
