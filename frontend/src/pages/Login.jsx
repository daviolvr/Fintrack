import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/Login.css";
import background from "../assets/images/register-login-background.jpg";
import logo from "../assets/images/logo_fintrack.png";
import { loginUser } from "../api/user_api";

export default function Login() {
    const navigate = useNavigate();

    // estados para os inputs
    const [form, setForm] = useState({
        email: "",
        password: "",
    });

    // estado para loading e mensagens
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");

    // atualizar os valores dos inputs
    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value});
        // limpa mensagem quando o usuario digitar
        if (message) setMessage("");
    };

    // validar formulario
    const validateForm = () => {
        if (!form.email || !form.password) {
            setMessage("Todos os campos são obrigatórios");
            return false;
        }

        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(form.email)) {
            setMessage("Por favor, insira um emáil válido");
            return false;
        }

        return true;
    };

    // submit do formulario
    const handleSubmit = async (e) => {
        e.preventDefault();
        
        if (!validateForm()) return;

        setLoading(true);
        setMessage("");

        try {
            // Envia os dados pro backend
            const response = await loginUser({
                email: form.email,
                password: form.password
            });

            setMessage("Login realizado com sucesso!");

            // Redireciona para a página de dashboard
            navigate("/dashboard");
        } catch(error) {
            setMessage(error.message || "Erro ao fazer login");
            setForm({
                email: "",
                password: "",
            });
        } finally {
            setLoading(false);
        }
    };

    // função para redirecionar para register
    const handleRegisterClick = () => {
        navigate("/register");
    }

    return (
        <div
            className="login-container"
            style={{ backgroundImage: `url(${background})` }}
        >
            <form onSubmit={handleSubmit} className="login-form">
                <img src={logo} alt="fintrack-logo" className="logo-image" />

                {/* <h2>Login</h2> */}

                {message && (
                    <div className={`message ${message.includes("sucesso") ? "success" : "error"}`}>
                        {message}
                    </div>
                )}

                <div className="input-group">
                    <label htmlFor="email">Email</label>
                    <input 
                        type="email" 
                        name="email" 
                        id="email" 
                        value={form.email} 
                        onChange={handleChange}
                        disabled={loading}
                        placeholder="Seu email"
                    />
                </div>

                <div className="input-group">
                    <label htmlFor="password">Senha</label>
                    <input 
                        type="password" 
                        name="password" 
                        id="password" 
                        value={form.password} 
                        onChange={handleChange}
                        disabled={loading}
                        placeholder="Sua senha"
                    />
                </div>

                <button className="login-button" type="submit" disabled={loading}>
                    {loading ? "Fazendo login..." : "Login"}
                </button>

                <button type="button" className="register-button" onClick={handleRegisterClick}>Criar uma conta</button>
            </form>
        </div>
    );
}