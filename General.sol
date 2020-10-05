pragma solidity ^0.7.2;

import "./AuthContract.sol";

contract GeneralContract {

    /* STRUCTS */
    struct User {
            string id;
            bool isRegistered;
            bool isLoggedIn;
            bool isAdmin;
            AuthContract auth;
            int8 attempts;
    }
    /* EVENTS */
    event createAdmin(address _adm, address _who);

    /* MODIFIERS */
    modifier isUser() {
        require(!userList[msg.sender].isRegistered, "This user is not in the system.");
        _;
    }

    modifier isAdmin() {
        require(!userList[msg.sender].isRegistered, "This user is not in the system.");
        require(userList[msg.sender].isAdmin, "This user does not have admin. priviledges.");
        _;
    }

    modifier userNotLocked{
        require (userList[msg.sender].attempts < 3);
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
        userList[owner].isRegistered = false;
        userList[owner].isAdmin = true;
    }

    /* -- ADMIN FUNCTIONS -- */

    function rmUser(address _addr, string memory _id) public isAdmin {
        // _addr = id2a[_id]
        userList[_addr].auth.terminate();
        userList[_addr].isRegistered = false;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
        id2a[_id] = address(0);
        userList[_addr].id = "";
    }

    function addUser(address _addr, string memory _id) public isAdmin {
        userList[_addr].auth = AuthContract(_addr);
        userList[_addr].isRegistered = true;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
        userList[_addr].id = _id;
        id2a[_id] = _addr;
    }

    function addAdmin(address _addr) public isAdmin {
        // Check that the user is added
        require(userList[_addr].isRegistered == false,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].isAdmin = true;
        // We notify in the blockchain who did it
        emit createAdmin(_addr, msg.sender);
    }

    /* -- USER FUNCTIONS (WRAPPERS)-- */
    function getOTP() public isUser userNotLocked returns(uint16 pass_){
        require (userList[msg.sender].isLoggedIn == false, "Only offline users can ask for a key");
        // We call that specific contract function
        pass_ = userList[msg.sender].auth.newOTP();
    }

    function tryLogin(uint16 _pass) public isUser userNotLocked {
        // We call that specific contract function
        require(userList[msg.sender].isLoggedIn == false, "You are already logged in");
        
        // Every time an attempt is made, the count is increased
        try userList[msg.sender].auth.tryLogin(_pass){
            // Successful login
            userList[msg.sender].attempts = 0;
        }
        catch {
            // Failed attempt
            userList[msg.sender].attempts++;
        }
    }
    
    function amILocked() public view isUser returns (bool locked){
        locked = (userList[msg.sender].attempts > 3);
    }

    function tryLogout() public isUser {
        // Can only logout if logged in
        require(userList[msg.sender].isLoggedIn == true, "You are not logged in");
        userList[msg.sender].isLoggedIn = false;
    }
}
