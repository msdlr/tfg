pragma solidity >=0.4.22 <0.7.0;

import "./AuthContract.sol";

contract GeneralContract {

    /* STRUCTS */
    struct User {
            string id;
            bool isNull;
            bool isLoggedIn;
            bool isAdmin;
            AuthContract auth;
    }
    /* EVENTS */
    event createAdmin(address _adm, address _who);

    /* MODIFIERS */
    modifier isUser() {
        require(!userList[msg.sender].isNull, "This user is not in the system.");
        _;
    }

    modifier isAdmin() {
        require(!userList[msg.sender].isNull, "This user is not in the system.");
        require(userList[msg.sender].isAdmin, "This user does not have admin. priviledges.");
        _;
    }

    /* Contract data */
    mapping ( address => User) userList;
    mapping ( string => address) id2a; // Index by id
    address owner;

    /* CONSTRUCTOR */

    constructor() public payable{
        // Set the owner of the company
        owner = msg.sender;

        // Add it to the admin list
        userList[owner].isNull = false;
        userList[owner].isAdmin = true;
    }

    /* -- ADMIN FUNCTIONS -- */

    function rmUser(address _addr, string _id) public isAdmin {
        // _addr = id2a[_id]
        userList[_addr].auth.terminate();
        userList[_addr].isNull = false;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
        id2a[_id] = address(0);
        userList[_addr].id = "";
    }

    function addUser(address _addr, string _id) public isAdmin {
        userList[_addr].auth = AuthContract(_addr);
        userList[_addr].isNull = true;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
        userList[_addr].id = _id;
        id2a[_id] = _addr;
    }

    function addAdmin(address _addr) public isAdmin {
        // Check that the user is added
        require(userList[_addr].isNull == false,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].isAdmin = true;
        // We notify in the blockchain who did it
        emit createAdmin(_addr, msg.sender);
    }

    /* -- USER FUNCTIONS (WRAPPERS)-- */
    function getOTP() public isUser returns(uint16 pass_){
        // We call that specific contract function
        pass_ = userList[msg.sender].Contract.newOTP();
    }

    function tryLogin(uint16 _pass) public isUser {
        // We call that specific contract function
        require(userList[_addr].isLoggedIn == false, "You are already logged in");
        userList[msg.sender].Contract.tryLogin(_pass);
    }

    function tryLogout() public isUser {
        // Can only logout if logged in
        require(userList[_addr].isLoggedIn == true, "You are not logged in");
        userList[_addr].auth.logout();
    }
}