import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie'

const AdminUsers = () => {
	const [users, setUsers] = useState([]);
	const [error, setError] = useState('');
	const navigate = useNavigate();

	useEffect(() => {
		const fetchUsers = async () => {
			try {
				const token = Cookies.get('token'); 
				if (!token) {
					setError('Необходима авторизация');
					return;
				}

				const response = await fetch('http://localhost:8080/admin/users', {
					method: 'GET',
					headers: {
						Authorization: `Bearer ${token}`,
					},
				});

				if (response.ok) {
					const data = await response.json();
					setUsers(data);
				} else if (response.status === 403) {
					setError(
						'Доступ запрещён: Только администратор может просматривать пользователей.'
					);
				} else {
					setError('Произошла ошибка при получении данных пользователей.');
				}
			} catch (err) {
				setError('Не удалось подключиться к серверу.');
			}
		};

		fetchUsers();
	}, []);

	return (
		<div>
			<h1>Панель администратора - Список пользователей</h1>
			{error && <p style={{ color: 'red' }}>{error}</p>}
			{!error && users.length > 0 && (
				<table border="1" cellPadding="10" cellSpacing="0">
					<thead>
						<tr>
							<th>ID</th>
							<th>Имя</th>
							<th>Email</th>
							<th>Телефон</th>
							<th>Роль</th>
							<th>Действия</th>
						</tr>
					</thead>
					<tbody>
						{users.map((user) => (
							<tr key={user.id}>
								<td>{user.ID}</td>
								<td>{user.Name}</td>
								<td>{user.Email}</td>
								<td>{user.Phone}</td>
								<td>{user.Role}</td>
								<td>
									<button
										onClick={() =>
											navigate(`/admin/users/${user.ID}/role`, {
												state: { user },
											})
										}
									>
										Изменить роль
									</button>
								</td>
							</tr>
						))}
					</tbody>
				</table>
			)}
			{!error && users.length === 0 && <p>Список пользователей пуст.</p>}
		</div>
	);
};

export default AdminUsers;
