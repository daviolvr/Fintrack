import React, { useState } from "react";
import "../styles/Register.css";
import background from "../assets/images/register-login-background.jpg";

export default function Register() {
    // estados para os inputs
    const [form, setForm] = useState({
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        confirmPassword: "",
    })

    // atualizar os valores dos inputs
    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value});
    };

    // submit do formulario
    const handleSubmit = (e) => {
        e.preventDefault();
        console.log("Dados do cadastro:", form);
    };

    return (
        <div
            className="register-container"
            style={{ backgroundImage: `url(${background})` }}
        >
            <form onSubmit={handleSubmit} className="register-form">
                <h2>Cadastro</h2>

                <div className="name-row">
                    <div className="input-group">
                        <label htmlFor="firstName">Nome</label>
                        <input type="text" name="firstName" id="firstName" value={form.firstName} onChange={handleChange} />
                    </div>
                    <div className="input-group">
                        <label htmlFor="lastName">Sobrenome</label>
                        <input type="text" name="lastName" id="lastName" value={form.lastName} onChange={handleChange} />
                    </div>
                </div>

                <div className="input-group">
                    <label htmlFor="email">Email</label>
                    <input type="email" name="email" id="email" value={form.email} onChange={handleChange} />
                </div>

                <div className="input-group">
                    <label htmlFor="password">Senha</label>
                    <input type="password" name="password" id="password" value={form.password} onChange={handleChange} />
                </div>

                <div className="input-group">
                    <label htmlFor="confirmPassword">Confirmar senha</label>
                    <input type="password" name="confirmPassword" id="confirmPassword" value={form.confirmPassword} onChange={handleChange} />
                </div>

                <button type="submit">Cadastrar</button>
            </form>
        </div>
    );
}