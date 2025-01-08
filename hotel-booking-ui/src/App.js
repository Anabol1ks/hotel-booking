import React from 'react'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import ResetPassword from './pages/ResetPassword'
import Login from './pages/Login'
import Register from './pages/Register'
import Home from './pages/Home'
import AdminUsers from './pages/AdminUsers'
import UpdateUserRole from './pages/UpdateUserRole'
import ForgotPassword from './pages/ForgotPassword'
import HotelsAndRooms from './pages/HotelsAndRooms'
import CreateOfflineBooking from './pages/CreateOfflineBooking'






const App = () => {
	return (
		<Router>
			<Routes>
				<Route path='/' element={<Home />} />
				<Route path='/auth/register' element={<Register />} />
				<Route path='/auth/login' element={<Login />} />
				<Route path='/auth/reset-password' element={<ResetPassword />} />
				<Route path='/auth/forgot-password' element={<ForgotPassword />} />
				<Route path='/admin/users' element={<AdminUsers />} />
				<Route path='/admin/users/:id/role' element={<UpdateUserRole />} />
				<Route path='/hotels-and-rooms' element={<HotelsAndRooms />} />
				<Route
					path='/bookings/offline/create'
					element={<CreateOfflineBooking />}
				/>
			</Routes>
		</Router>
	)
}


export default App
