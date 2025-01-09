import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './Home.module.css'

const Home = () => {
	const [role, setRole] = useState(null)
	const navigate = useNavigate()
	console.log(process.env.REACT_APP_API_URL)
	useEffect(() => {
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
		<div className={styles.container}>
			<h1 className={styles.title}>Главная страница</h1>
			<div className={styles.buttonContainer}>
				<button 
					className={`${styles.button} ${styles.primaryButton}`}
					onClick={() => navigate('/hotels-and-rooms')}
				>
					Просмотр отелей и номеров
				</button>
				
				{(role === 'manager' || role === 'owner') && (
					<button 
						className={`${styles.button} ${styles.primaryButton}`}
						onClick={() => navigate('/bookings/offline/create')}
					>
						Создать офлайн бронирование
					</button>
				)}

				{role === 'admin' && (
					<button 
						className={`${styles.button} ${styles.adminButton}`}
						onClick={() => navigate('/admin/users')}
					>
						Перейти в панель администратора
					</button>
				)}
				{role && (
					<button 
						className={`${styles.button} ${styles.secondaryButton}`}
						onClick={handleLogout}
					>
						Выйти
					</button>
				)}
				{!role && (
					<>
						<button 
							className={`${styles.button} ${styles.primaryButton}`}
							onClick={() => navigate('auth/login')}
						>
							Войти
						</button>
						<button 
							className={`${styles.button} ${styles.primaryButton}`}
							onClick={() => navigate('auth/register')}
						>
							Зарегистрироваться
						</button>
					</>
				)}
			</div>
		</div>
	)
}

export default Home