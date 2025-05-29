import React, { useState, useEffect } from "react";
import axios from "axios";

export default function Auth({
  setIsAuthenticated,
  fetchMessages,
  isAdmin,
  setIsAdmin,
}) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoginDisabled, setIsLoginDisabled] = useState(false);
  const [blockTimeRemaining, setBlockTimeRemaining] = useState(0);
  const [remainingAttempts, setRemainingAttempts] = useState(null);

  const MAX_ATTEMPTS = 3;
  const BLOCK_TIME_MS = 30000; // 30 секунд блокировки

  useEffect(() => {
    const blockUntil = localStorage.getItem("loginBlockUntil");
    if (blockUntil && Date.now() < parseInt(blockUntil)) {
      const remainingTime = Math.ceil(
        (parseInt(blockUntil) - Date.now()) / 1000
      );
      setBlockTimeRemaining(remainingTime);
      setIsLoginDisabled(true);
      setError(
        `Слишком много попыток. Попробуйте через ${remainingTime} секунд.`
      );

      const timer = setInterval(() => {
        const newRemaining = Math.ceil(
          (parseInt(blockUntil) - Date.now()) / 1000
        );
        if (newRemaining <= 0) {
          setIsLoginDisabled(false);
          setError("");
          localStorage.removeItem("loginBlockUntil");
          clearInterval(timer);
        } else {
          setBlockTimeRemaining(newRemaining);
          setError(
            `Слишком много попыток. Попробуйте через ${newRemaining} секунд.`
          );
        }
      }, 1000);

      return () => clearInterval(timer);
    }
  }, []);

  const handleLogin = async (e) => {
    e.preventDefault();

    if (isLoginDisabled) {
      setError(
        `Слишком много попыток. Попробуйте через ${blockTimeRemaining} секунд.`
      );
      return;
    }

    try {
      const response = await axios.post("http://localhost:5000/login", {
        login,
        password,
      });

      localStorage.setItem("token", response.data.token);
      localStorage.removeItem("failedLoginAttempts");
      localStorage.removeItem("loginBlockUntil");
      setIsAuthenticated(true);
      setError("");
      setRemainingAttempts(null);
      await fetchMessages();
    } catch (err) {
      const attempts =
        parseInt(localStorage.getItem("failedLoginAttempts") || 0) + 1;
      localStorage.setItem("failedLoginAttempts", attempts.toString());

      if (attempts >= MAX_ATTEMPTS) {
        const blockUntil = Date.now() + BLOCK_TIME_MS;
        localStorage.setItem("loginBlockUntil", blockUntil.toString());
        setIsLoginDisabled(true);
        setBlockTimeRemaining(Math.ceil(BLOCK_TIME_MS / 1000));
        setError(
          `Слишком много попыток. Попробуйте через ${Math.ceil(
            BLOCK_TIME_MS / 1000
          )} секунд.`
        );

        const timer = setInterval(() => {
          const newRemaining = Math.ceil((blockUntil - Date.now()) / 1000);
          if (newRemaining <= 0) {
            setIsLoginDisabled(false);
            setError("");
            clearInterval(timer);
          } else {
            setError(
              `Слишком много попыток. Попробуйте через ${newRemaining} секунд.`
            );
          }
        }, 1000);
      } else {
        const remaining = MAX_ATTEMPTS - attempts;
        setRemainingAttempts(remaining);
        setError(`Неверные учетные данные. Осталось попыток: ${remaining}`);
      }
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-6 text-center">Авторизация</h1>

        {error && (
          <div
            className={`mb-4 p-2 rounded text-center ${
              isLoginDisabled
                ? "bg-red-100 text-red-700"
                : "bg-yellow-100 text-yellow-700"
            }`}
          >
            {error}
          </div>
        )}

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Логин
            </label>
            <input
              type="text"
              value={login}
              onChange={(e) => setLogin(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              disabled={isLoginDisabled}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Пароль
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              disabled={isLoginDisabled}
            />
          </div>
          <button
            type="submit"
            disabled={isLoginDisabled}
            className={`w-full py-2 px-4 rounded-lg transition-colors ${
              isLoginDisabled
                ? "bg-gray-400 cursor-not-allowed"
                : "bg-blue-600 text-white hover:bg-blue-700"
            }`}
          >
            {isLoginDisabled ? "Вход временно заблокирован" : "Войти"}
          </button>
        </form>
      </div>
    </div>
  );
}
