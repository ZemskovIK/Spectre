import React, { useState, useEffect } from "react";
import axios from "axios";
import ModalWindow from "./ModalWindow";
import EmailList from "./EmailList";
import SendFormSection from "./SendFormSection";
import EditFormSection from "./EditFormSection";

export default function MailApp() {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [author, setAuthor] = useState("");
  const [foundIn, setFoundIn] = useState("");
  const [dataInt, setDataInt] = useState("");
  const [newEditingMessage, setEditingNewMessage] = useState("");
  const [editingAuthor, setEditingAuthor] = useState("");
  const [editingFoundIn, setEditingFoundIn] = useState("");
  const [editingDataInt, setEditingDataInt] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [deletingId, setDeletingId] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [isntError, setIsntError] = useState();
  const [errorText, setErrorText] = useState();

  useEffect(() => {
    fetchMessages();
  }, [messages]);

  const fetchMessages = async () => {
    try {
      const response = await axios.get("http://localhost:5000/api/letters");
      setMessages(response.data);
    } catch (error) {
      console.error("Ошибка отправки:", error);
    }
  };

  const deleteMessage = async (id) => {
    setDeletingId(id);
    try {
      await axios.delete(`http://localhost:5000/api/letters/${id}`);
    } catch (error) {
      console.error("ошибка удаления:", error);
    } finally {
      setMessages(messages.filter((msg) => msg.id !== id));
      setDeletingId(null);
    }
  };

  const editMessage = async (e) => {
    e.preventDefault();
    if (!newEditingMessage.trim() || !editingAuthor.trim()) return;

    setLoading(true);
    try {
      console.log(newEditingMessage);
      const response = await axios.put(
        `http://localhost:5000/api/letters/${editingId}`,
        {
          body: newEditingMessage,
          author: editingAuthor.trim(),
          found_at: editingDataInt.trim(),
          found_in: editingFoundIn.trim(),
          // id2: "1",
        }
      );
      setMessages([...messages, response.data]);
      setEditingNewMessage("");
      setEditingAuthor("");
      setEditingDataInt("");
      setEditingFoundIn("");
      setEditingId("");
    } catch (error) {
      console.error("ошибка отправки:", error);
    } finally {
      setLoading(false);
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !author.trim()) return;

    setLoading(true);
    try {
      const response = await axios.post("http://localhost:5000/api/letters", {
        body: newMessage,
        author: author.trim(),
        found_at: dataInt.trim(),
        found_in: foundIn.trim(),
        // id2: "1",
      });
      setMessages([...messages, response.data]);
      setNewMessage("");
      setAuthor("");
      setDataInt("");
      setFoundIn("");
    } catch (error) {
      console.error("ошибка отправки:", error);
    } finally {
      setLoading(false);
    }
  };

  const openMessage = (message, er, ertxt) => {
    setSelectedMessage(message);
    setIsntError(er);
    setErrorText(ertxt);
    console.log(er);
    setShowModal(true);
  };

  return (
    <div className="min-h-screen m-5">
      {/* <div className="max-w-6xl mx-auto mb-6">
        <h1 className="text-3xl font-bold text-blue-600">Spectre</h1>
        <p className="text-gray-600">какая нибудь шапка</p>
      </div> */}

      <EmailList
        messages={messages}
        onEmailClick={openMessage}
        onDeleteEmail={deleteMessage}
        deletingId={deletingId}
      ></EmailList>
      {/* w-full md:w-1/2 lg:w-3/4 xl:w-2/3 */}
      <div className="max-w-6xl mx-auto mt-6 bg-gradient-to-r from-gray-200 via-gray-350 to-gray-400 rounded-lg shadow-md p-6">
        <button
          onClick={() => setIsEditing(false)}
          className={`py-1 px-1 mr-3 rounded-lg text-white font-medium text-2xl font-semibold ${
            isEditing
              ? "bg-gray-500 hover:bg-gray-700"
              : "bg-gray-600 hover:bg-gray-700"
          }`}
        >
          Добавить новое письмо
        </button>{" "}
        <button
          onClick={() => setIsEditing(true)}
          className={`py-1 px-1 rounded-lg text-white font-medium text-2xl font-semibold mb-4 ${
            isEditing
              ? "bg-gray-600 hover:bg-gray-700"
              : "bg-gray-500 hover:bg-gray-700"
          }`}
        >
          Редактировать письмо
        </button>
        {isEditing == false ? (
          <>
            {" "}
            <form
              onSubmit={sendMessage}
              className="space-y-4"
              autocomplete="off"
            >
              <SendFormSection
                author={author}
                setAuthor={setAuthor}
                foundIn={foundIn}
                setFoundIn={setFoundIn}
                dataInt={dataInt}
                setDataInt={setDataInt}
                newMessage={newMessage}
                setNewMessage={setNewMessage}
              ></SendFormSection>

              <button
                type="submit"
                disabled={loading || !newMessage.trim() || !author.trim()}
                className={`px-6 py-2 rounded-lg text-white font-medium ${
                  loading || !newMessage.trim() || !author.trim()
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-gray-600 hover:bg-gray-700"
                }`}
              >
                {loading ? "отправка..." : "отправить"}
              </button>
            </form>
          </>
        ) : (
          <form onSubmit={editMessage} className="space-y-4" autocomplete="off">
            <EditFormSection
              editingId={editingId}
              setEditingId={setEditingId}
              editingAuthor={editingAuthor}
              editingDataInt={editingDataInt}
              editingFoundIn={editingFoundIn}
              newEditingMessage={newEditingMessage}
              setEditingAuthor={setEditingAuthor}
              setEditingDataInt={setEditingDataInt}
              setEditingFoundIn={setEditingFoundIn}
              setEditingNewMessage={setEditingNewMessage}
            ></EditFormSection>
            <button
              type="submit"
              disabled={
                loading || !newEditingMessage.trim() || !editingAuthor.trim()
              }
              className={`px-6 py-2 rounded-lg text-white font-medium ${
                loading || !newEditingMessage.trim() || !editingAuthor.trim()
                  ? "bg-gray-300 cursor-not-allowed"
                  : "bg-gray-600 hover:bg-gray-700"
              }`}
            >
              {loading ? "отправка..." : "отправить"}
            </button>
          </form>
        )}
      </div>

      {showModal && (
        <ModalWindow
          onClose={() => setShowModal(false)}
          isError={isntError}
          errorText={errorText}
          dataInt={selectedMessage.found_at}
          title={selectedMessage.author}
          foundIn={selectedMessage.found_in}
          text={selectedMessage.body}
          backgroundImage="https://static.vecteezy.com/system/resources/previews/032/048/239/non_2x/paper-vintage-background-recycle-brown-paper-crumpled-texture-ai-generated-free-photo.jpg"
        ></ModalWindow>
      )}
    </div>
  );
}
