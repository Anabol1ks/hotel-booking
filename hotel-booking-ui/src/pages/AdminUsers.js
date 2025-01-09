import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie'
import styles from './AdminUsers.module.css'

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

				const response = await fetch(process.env.REACT_APP_API_URL + '/admin/users', {
					method: 'GET',
					headers: {
						Authorization: `Bearer ${token}`,
					},
				});

				if (response.ok) {
					const data = await response.json();
					setUsers(data);
				} else if (response.status === 403) {
					setError('Доступ запрещён: Только администратор может просматривать пользователей.');
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
		<div className={styles.container}>
			<h1 className={styles.title}>
				Панель администратора - Список пользователей
			</h1>
			{error && <p className={styles.error}>{error}</p>}
			{!error && users.length > 0 && (
				<table className={styles.table}>
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
						{users.map(user => (
							<tr key={user.id}>
								<td>{user.ID}</td>
								<td>{user.Name}</td>
								<td>{user.Email}</td>
								<td>{user.Phone}</td>
								<td>{user.Role}</td>
								<td>
									<button
										className={styles.actionButton}
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
			<button className={styles.backButton} onClick={() => navigate('/')}>
				Вернуться на главную
			</button>
			<h1 className={styles.title}>
				Панель администратора - Список пользователей
			</h1>
			{/* Rest of the existing code */}
		</div>
	)
	
};

export default AdminUsers;