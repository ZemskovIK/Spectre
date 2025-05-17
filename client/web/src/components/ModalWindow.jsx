import React from "react";
import "./Modal.css";

export default function Modal({
  isError,
  errorText,
  title,
  text,
  onClose,
  backgroundImage,
  foundIn,
  dataInt,
}) {
  return (
    <>
      <div className="modal-blur-overlay" />

      <div className="modal-container">
        <div
          className="modal-background"
          style={{ backgroundImage: `url(${backgroundImage})` }}
        />

        <div className="modal-content">
          <button
            className="close-button hover:scale-120 transition-transform duration-200"
            onClick={onClose}
          >
            &#10006;
          </button>

          <div className="modal-header">
            {" "}
            <h2 className="modal-title">Автор: {title}</h2>
          </div>

          <div className="modal-body">
            {" "}
            <div className="">
              {" "}
              <h2 className="modal-found">Найдено: {foundIn}</h2>
            </div>
            <br></br>
            <div className="">
              {" "}
              <h2 className="modal-found">Дата: {dataInt}</h2>
            </div>
            <br></br>
            <p className="modal-text">{text}</p>
          </div>
        </div>
      </div>
    </>
  );
}
