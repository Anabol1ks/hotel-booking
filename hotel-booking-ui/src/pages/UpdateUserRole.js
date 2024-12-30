import React, { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

const UpdateUserRole = () => {
	const location = useLocation(); // Получение переданных данных через state
	const navigate = useNavigate();

	const user = location.state?.user; // Извлечение данных пользователя из state
	const [role, setRole] = useState(user?.role || '');
	const [error, setError] = useState('');
	const [success, setSuccess] = useState('');
	const [loading, setLoading] = useState(false);

	// Список доступных ролей
	const roles = ['client', 'owner', 'admin'];

	const handleSubmit = async (e) => {
		e.preventDefault();

		setLoading(true);
		setError('');
		setSuccess('');

		try {
			const token = localStorage.getItem('token');

			const response = await fetch(
				`http://localhost:8080/admin/users/${user.ID}/role`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${token}`,
					},
					body: JSON.stringify({ role }),
				}
			);

			if (response.ok) {
				setSuccess('Роль успешно обновлена');
			} else {
				const data = await response.json();
				setError(data.error || 'Ошибка при обновлении роли');
			}
		} catch (err) {
			setError('Произошла ошибка. Попробуйте еще раз.');
		} finally {
			setLoading(false);
		}
	};

	if (!user) {
		return <p>Данные о пользователе не переданы</p>;
	}

	return (
		<div>
			<h1>Смена роли пользователя</h1>
			<p>Редактирование пользователя: {user.Name}</p>
			<p>Email: {user.Email}</p>

			{error && <p style={{ color: 'red' }}>{error}</p>}
			{success && <p style={{ color: 'green' }}>{success}</p>}

			<form onSubmit={handleSubmit}>
				<div>
					<label htmlFor="role">Выберите новую роль:</label>
					<select
						id="role"
						value={role}
						onChange={(e) => setRole(e.target.value)}
						required
					>
						<option value="" disabled>
							Выберите роль
						</option>
						{roles.map((r) => (
							<option key={r} value={r}>
								{r}
							</option>
						))}
					</select>
				</div>
				<button type="submit" disabled={loading}>
					{loading ? 'Обновление...' : 'Обновить роль'}
				</button>
			</form>
			<button onClick={() => navigate('/admin/users')}>Вернуться назад</button>
		</div>
	);
};

export default UpdateUserRole;
