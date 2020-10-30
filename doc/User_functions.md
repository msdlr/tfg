### User functions.

User functions are accessed from the appication's general contract, although being in other separate contects like AuthContract or such, through wrapper functions; allowing the access and keeping control of these contracts from the main one.

#### Login and authentication.
Users are added by their Ethereum address to the system by administrators, making the management of users and contract centralised within the General Contract, but still retaining distributed access to that contract. 
Users have 2 functions related to login and authentication:

 - Get OTP (one time password): the function can only be called by users that are registered in the system, and that are not locked out from it. This generates a numeric password from 1 to 9999 for the user to input when logging in, which lasts for approximately . This temporary password is both returned from the function and stored in the blockchain hashed using the kekkak/sha3 algorithm.

 - Try login: this function has as a pre-requisite that the user has requested a valid OTP that has not expired. If the number passed to this function is hashed and results in the hash stored in the blockchain, meaning that it is correct, the General Contract will take not of the user being connected. There's also an attempts count per user, which is increased upon failed login attempts, and reset upon success. When this counter gets to 3, the user is locked from the system, which would require an ademinitstrator's intervention to reset their attempt counter. 
