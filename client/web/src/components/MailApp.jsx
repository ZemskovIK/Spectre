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
  const [isError, setIsError] = useState("Писем нет");
  const [searchQuery, setSearchQuery] = useState("");
  // const [listOfIndex, setListOfIndex] = useState([]);

  useEffect(() => {
    fetchMessages();

    const intervalId = setInterval(fetchMessages, 60000);

    return () => clearInterval(intervalId);
  }, []);

  const fetchMessages = async () => {
    try {
      const response = await axios.get("http://localhost:5000/api/letters");
      if (response.data.error == null) {
        setMessages(response.data.content);
        // setListOfIndex(response.data.content.map((item) => item.id));
      } else {
        setIsError(response.data.error);
      }
    } catch (error) {
      console.error("ошибка отправки:", error);
    }
  };

  const filteredMessages = searchQuery.trim()
    ? messages.filter((message) =>
        message?.author?.toLowerCase()?.includes(searchQuery.toLowerCase())
      )
    : messages;

  function transDataToDima(data) {
    let newVal = `${data.slice(6, 10)}-${data.slice(3, 5)}-${data.slice(
      0,
      2
    )}T00:00:00Z`;
    return newVal;
  }

  function handleChangeModes() {
    setIsEditing(!isEditing);
    setEditingDataInt("");
    setDataInt("");
  }

  const handleChange = (e) => {
    let value = e.target.value;

    value = value.replace(/[^\d]/g, "");

    if (value.length > 2) {
      value = `${value.slice(0, 2)}.${value.slice(2)}`;
    }
    if (value.length > 5) {
      value = `${value.slice(0, 5)}.${value.slice(5, 9)}`;
    }

    if (value.length > 10) return;

    setDataInt(value);
    setEditingDataInt(value);
  };

  const deleteMessage = async (id) => {
    setDeletingId(id);
    try {
      await axios.delete(`http://localhost:5000/api/letters/${id}`);
      fetchMessages();
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
          found_at: transDataToDima(editingDataInt.trim()),
          found_in: editingFoundIn.trim(),
          // id2: "1",
        }
      );
      fetchMessages();
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
        found_at: transDataToDima(dataInt.trim()),
        found_in: foundIn.trim(),
        // id2: "1",
      });
      fetchMessages();
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

  const openMessage = (message) => {
    setSelectedMessage(message);
    setShowModal(true);
  };

  return (
    <div className="min-h-screen m-5">
      {/* <div className="max-w-6xl mx-auto mb-6">
        <h1 className="text-3xl font-bold text-blue-600">Spectre</h1>
        <p className="text-gray-600">какая нибудь шапка</p>
      </div> */}

      <EmailList
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        isError={isError}
        messages={filteredMessages}
        onEmailClick={openMessage}
        onDeleteEmail={deleteMessage}
        deletingId={deletingId}
      ></EmailList>
      {/* w-full md:w-1/2 lg:w-3/4 xl:w-2/3 */}
      <div className="lg:w-6xl 2xl:w-7xl mx-auto mt-6 bg-gradient-to-r from-gray-200 via-gray-350 to-gray-400 rounded-lg shadow-md p-6">
        <button
          onClick={handleChangeModes}
          className={`py-1 px-1 mr-3 rounded-lg text-white font-medium text-2xl font-semibold ${
            isEditing
              ? "bg-gray-500 hover:bg-gray-700"
              : "bg-gray-600 hover:bg-gray-700"
          }`}
        >
          Добавить новое письмо
        </button>{" "}
        <button
          onClick={handleChangeModes}
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
              className="space-y-2"
              autocomplete="off"
            >
              <SendFormSection
                handleChange={handleChange}
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
                disabled={
                  loading ||
                  !newMessage.trim() ||
                  !author.trim() ||
                  dataInt.length < 10
                }
                className={`px-6 py-2 rounded-lg text-white font-medium ${
                  loading ||
                  !newMessage.trim() ||
                  !author.trim() ||
                  dataInt.length < 10
                    ? "bg-gray-300 cursor-not-allowed"
                    : "bg-gray-600 hover:bg-gray-700"
                }`}
              >
                {loading ? "отправка..." : "отправить"}
              </button>
            </form>
          </>
        ) : (
          <form onSubmit={editMessage} className="space-y-2" autocomplete="off">
            <EditFormSection
              handleChange={handleChange}
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
                loading ||
                !newEditingMessage.trim() ||
                !editingAuthor.trim() ||
                editingDataInt.length < 10
              }
              className={`px-6 py-2 rounded-lg text-white font-medium ${
                loading ||
                !newEditingMessage.trim() ||
                !editingAuthor.trim() ||
                editingDataInt.length < 10
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
