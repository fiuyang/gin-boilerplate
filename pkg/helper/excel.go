package helper

var ColumnCustomer = map[int]string{
	0: "username is required",
	1: "email is required",
	2: "phone is required",
	3: "address is required",
}

var UniqueCustomer = map[int]string{
	0: "email",
}

var RulesUser = map[int]string{
	0: "username,required",
	1: "email,required,unique",
	2: "password,required",
}
