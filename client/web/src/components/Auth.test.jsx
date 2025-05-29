import React from "react";
import { render, screen } from "@testing-library/react";
import Auth from "./Auth";

describe("Auth Component", () => {
  const mockProps = {
    setIsAuthenticated: jest.fn(),
    fetchMessages: jest.fn(),
    isAdmin: null,
    setIsAdmin: jest.fn(),
  };

  it("renders login form correctly", () => {
    render(<Auth {...mockProps} />);

    // Проверяем заголовок
    expect(
      screen.getByRole("heading", { name: /авторизация/i })
    ).toBeInTheDocument();

    // Проверяем поля ввода по их ролям
    expect(screen.getByRole("textbox", { name: /логин/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/пароль/i)).toBeInTheDocument();

    // Проверяем кнопку
    expect(screen.getByRole("button", { name: /войти/i })).toBeInTheDocument();
  });
});
