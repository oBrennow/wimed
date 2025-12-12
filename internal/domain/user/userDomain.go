package user


type User struct{
	id		string
	email	string
	roles	[]roles
	active	bool
}

type roles struct {
	nivel	string
	acessos	[]acessos
}


func CreatNewUserDomain(id string, email string, roles []roles, active bool)(*User, err error){
	return nil, err
}