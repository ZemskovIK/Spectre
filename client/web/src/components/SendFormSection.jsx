export default function SendFormSection({
  handleChange,
  author,
  setAuthor,
  foundIn,
  setFoundIn,
  dataInt,
  setDataInt,
  newMessage,
  setNewMessage,
}) {
  return (
    <>
      <div className="flex space-x-4">
        <div className="flex flex-col w-full">
          <label
            htmlFor="author"
            className=" text-sm font-medium text-gray-700 mb-1"
          >
            Автор
          </label>
          <input
            id="author"
            type="text"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            placeholder="автор"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col">
          <label
            htmlFor="найдено"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            найдено
          </label>
          <input
            id="author"
            type="text"
            value={foundIn}
            onChange={(e) => setFoundIn(e.target.value)}
            placeholder="найдено"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col">
          <label
            htmlFor="дата"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            дата
          </label>
          <input
            id="date"
            type="text"
            value={dataInt}
            onChange={handleChange}
            placeholder="дата дд.мм.гггг"
            className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
      </div>
      <div>
        <label
          htmlFor="письмо"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Письмо
        </label>
        <textarea
          id="message"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="пишите тут"
          className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          required
        />
      </div>
    </>
  );
}
