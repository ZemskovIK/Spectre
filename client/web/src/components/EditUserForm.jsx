export default function EditUserForm({
  editingUserId,
  setEditingUserId,
  editingUserLogin,
  setEditingUserLogin,
  editingUserPassword,
  setEditingUserPassword,
  editingUserAccessLevel,
  setEditingUserAccessLevel,
}) {
  return (
    <>
      <div className="flex space-x-4">
        <div className="flex flex-col w-full">
          <label
            htmlFor="id"
            className=" text-sm font-medium text-gray-700 mb-1"
          >
            Id редактируемого пользователя:
          </label>
          <input
            id="aid"
            type="text"
            value={editingUserId}
            onChange={(e) => setEditingUserId(e.target.value)}
            placeholder="Id редактируемого пользователя"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col w-full">
          <label
            htmlFor="login"
            className=" text-sm font-medium text-gray-700 mb-1"
          >
            Логин
          </label>
          <input
            id="login"
            type="text"
            value={editingUserLogin}
            onChange={(e) => setEditingUserLogin(e.target.value)}
            placeholder="новый логин"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col w-full">
          <label
            htmlFor="accesslvl"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Уровень доступа
          </label>
          <input
            id="author"
            type="text"
            value={editingUserAccessLevel}
            onChange={(e) => setEditingUserAccessLevel(e.target.value)}
            placeholder="новый уровень доступа"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
      </div>

      <div>
        <label
          htmlFor="Пароль"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Пароль
        </label>
        <textarea
          id="message"
          value={editingUserPassword}
          onChange={(e) => setEditingUserPassword(e.target.value)}
          placeholder="новый пароль"
          className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          required
        />
      </div>
    </>
  );
}
