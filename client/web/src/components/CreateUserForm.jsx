export default function CreateUserForm({
  setUserId,
  userLogin,
  setUserLogin,
  userPassword,
  setUserPassword,
  userAccessLevel,
  setUserAccessLevel,
}) {
  return (
    <>
      <div className="flex space-x-4">
        <div className="flex flex-col w-full">
          <label
            htmlFor="author"
            className=" text-sm font-medium text-gray-700 mb-1"
          >
            Логин
          </label>
          <input
            id="login"
            type="text"
            value={userLogin}
            onChange={(e) => setUserLogin(e.target.value)}
            placeholder="логин"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="flex flex-col">
          <label
            htmlFor="найдено"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Уровень доступа
          </label>
          <input
            id="accesslvl"
            type="text"
            value={userAccessLevel}
            onChange={(e) => setUserAccessLevel(e.target.value)}
            placeholder="уровень доступа"
            className="block px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
      </div>
      <div>
        <label
          htmlFor="письмо"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Пароль
        </label>
        <textarea
          id="message"
          value={userPassword}
          onChange={(e) => setUserPassword(e.target.value)}
          placeholder="пароль"
          className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows={4}
          required
        />
      </div>
    </>
  );
}
