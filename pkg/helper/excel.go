package helper

var ColumnExcelCustomer = map[int]string{
	0: "username is required",
	1: "email is required",
	2: "phone is required",
	3: "address is required",
}

var UniqueExcelCustomer = map[int]string{
	0: "email",
}

var RulesExcelUser = map[int]string{
	0: "username,required",
	1: "email,required,unique",
	2: "password,required",
}
