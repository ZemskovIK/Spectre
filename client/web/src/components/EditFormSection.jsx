export default function EditFormSection({
  editingId,
  setEditingId,
  editingAuthor,
  setEditingAuthor,
  editingDataInt,
  setEditingDataInt,
  editingFoundIn,
  setEditingFoundIn,
  newEditingMessage,
  setEditingNewMessage,
}) {
  return (
    <>
      <div className="flex space-x-4">
        <div className="flex flex-col w-full">
          <label
            htmlFor="id"
            className=" text-sm font-medium text-gray-700 mb-1"
          >
            Номер редактируемого письма:
          </label>
          <input
            id="author"
            type="text"
            value={editingId}
            onChange={(e) => setEditingId(e.target.value)}
            placeholder="Номер редактируемого письма"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
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
            value={editingAuthor}
            onChange={(e) => setEditingAuthor(e.target.value)}
            placeholder="автор"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col w-full">
          <label
            htmlFor="найдено"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Найдено
          </label>
          <input
            id="author"
            type="text"
            value={editingFoundIn}
            onChange={(e) => setEditingFoundIn(e.target.value)}
            placeholder="найдено"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col w-full">
          <label
            htmlFor="дата"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Дата
          </label>
          <input
            id="author"
            type="text"
            value={editingDataInt}
            onChange={(e) => setEditingDataInt(e.target.value)}
            placeholder="дата"
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
          value={newEditingMessage}
          onChange={(e) => setEditingNewMessage(e.target.value)}
          placeholder="пишите тут"
          className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          required
        />
      </div>
    </>
  );
}
