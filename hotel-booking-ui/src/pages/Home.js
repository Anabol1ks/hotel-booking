import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'

const Home = () => {
	const [role, setRole] = useState(null)
	const navigate = useNavigate()

	useEffect(() => {
		// Получаем роль из localStorage
		const storedRole = Cookies.get('role')
		if (storedRole) {
			setRole(storedRole)
		}
	}, [])

	const handleLogout = () => {
		Cookies.remove('token')
		Cookies.remove('role')
		navigate('auth/login')
	}

	return (
		<div>
			<h1>Главная страница</h1>
			{role === 'admin' && (
				<button onClick={() => navigate('/admin/users')}>
					Перейти в панель администратора
				</button>
			)}
			{role && <button onClick={handleLogout}>Выйти</button>}
			{!role && (
				<div>
					<button onClick={() => navigate('auth/login')}>Войти</button>
					<button onClick={() => navigate('auth/register')}>
						Зарегистрироваться
					</button>
				</div>
			)}
		</div>
	)
}

export default Home
