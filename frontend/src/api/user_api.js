const API_BASE_URL = "http://localhost:8000/api/v1";

// Função para registrar usuário
export const registerUser = async (userData) => {
    try {
        const response = await fetch(`${API_BASE_URL}/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(userData),
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || "Erro ao cadastrar usuário");
        }

        return data;
    } catch (error) {
        throw error;
    }
};

// Função para fazer login
export const loginUser = async(userData) => {
    try {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(userData)
        });

        const data = await response.json();

        if(!response.ok) {
            throw new Error(data.message || "Erro ao fazer login");
        }

        return data;
    } catch (error) {
        throw error;
    }
}
