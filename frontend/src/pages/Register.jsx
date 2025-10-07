import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/Register.css";
import background from "../assets/images/register-login-background.jpg";
import logo from "../assets/images/logo_fintrack.png";
import { registerUser } from "../api/user_api";

const MESSAGES = {
    REQUIRED_FIELDS: "Todos os campos são obrigatórios",
    INVALID_EMAIL: "Por favor, insira um email válido",
    MISMATCHING_PASSWORDS: "As senhas não coincidem",
    INVALID_PASSWORD: "A senha deve ter pelo menos 8 caracteres",
    REGISTER_SUCCESS: "Cadastro realizado com sucesso!",
    REGISTER_ERROR: "Erro ao cadastrar usuário"
};

export default function Register() {
    const navigate = useNavigate();

    // estados para os inputs
    const [form, setForm] = useState({
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        confirmPassword: "",
    });

    // estado para loading e mensagens
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");

    // atualizar os valores dos inputs
    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value});
        // limpa mensagem quando o usuário digitar
        if (message) setMessage("");
    };

    // validar email com domínios permitidos
    const validateEmailDomain = (email) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            return false;
        }

        // Domínios permitidos
        const allowedDomains = [
            'gmail.com',
            'outlook.com',
        ];

        const domain = email.split('@')[1].toLowerCase();
        return allowedDomains.includes(domain);
    };

    // validar formulário
    const validateForm = () => {
        if (!form.firstName || !form.lastName || !form.email || !form.password) {
            setMessage(MESSAGES.REQUIRED_FIELDS);
            return false;
        }

        if (form.password !== form.confirmPassword) {
            setMessage(MESSAGES.MISMATCHING_PASSWORDS);
            return false;
        }

        if (form.password.length < 8) {
            setMessage(MESSAGES.INVALID_PASSWORD);
            return false;
        }

        if (!validateEmailDomain(form.email)) {
            setMessage(MESSAGES.INVALID_EMAIL);
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
            // Montar objeto para a API
            const userData = {
                first_name: form.firstName,
                last_name: form.lastName,
                email: form.email,
                password: form.password,
            };

            await registerUser(userData);

            setMessage(MESSAGES.REGISTER_SUCCESS);
            
            setTimeout(() => {
                navigate("/dashboard");
            }, 2000);
        } catch (error) {
            setMessage(error.message || MESSAGES.REGISTER_ERROR);
            setForm({
                firstName: "",
                lastName: "",
                email: "",
                password: "",
                confirmPassword: "",
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div
            className="register-container"
            style={{ backgroundImage: `url(${background})` }}
        >
            <form onSubmit={handleSubmit} className="register-form">
                <img src={logo} alt="fintrack-logo" className="logo-image" />

                {/* <h2>Cadastro</h2> */}

                {message && (
                    <div className={`message ${message.includes("sucesso") ? "success" : "error"}`}>
                        {message}
                    </div>
                )}

                <div className="name-row">
                    <div className="input-group">
                        <label htmlFor="firstName">Nome</label>
                        <input 
                            type="text" 
                            name="firstName" 
                            id="firstName" 
                            value={form.firstName} 
                            onChange={handleChange}
                            disabled={loading}
                            placeholder="Seu nome"
                        />
                    </div>
                    <div className="input-group">
                        <label htmlFor="lastName">Sobrenome</label>
                        <input 
                            type="text" 
                            name="lastName" 
                            id="lastName" 
                            value={form.lastName} 
                            onChange={handleChange}
                            disabled={loading}
                            placeholder="Seu sobrenome"
                        />
                    </div>
                </div>

                <div className="input-group">
                    <label htmlFor="email">Email</label>
                    <input 
                        type="email" 
                        name="email" 
                        id="email" 
                        value={form.email} 
                        onChange={handleChange}
                        disabled={loading}
                        placeholder="exemplo@gmail.com"
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
                        placeholder="Mínimo 8 caracteres"
                    />
                </div>

                <div className="input-group">
                    <label htmlFor="confirmPassword">Confirmar senha</label>
                    <input 
                        type="password" 
                        name="confirmPassword" 
                        id="confirmPassword" 
                        value={form.confirmPassword} 
                        onChange={handleChange}
                        disabled={loading}
                        placeholder="Digite a senha novamente"
                    />
                </div>

                <button type="submit" disabled={loading}>
                    {loading ? "Cadastrando..." : "Cadastrar"}
                </button>
            </form>
        </div>
    );
}