pragma solidity >=0.4.22 <0.7.0;

import "./Auth.sol";
import "./_structs_events.sol";

contract GeneralContract {

    /* STRUCTS */

    /* EVENTS */
    event createAdmin(address _adm, address _who, uint256 t);

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
    //mapping ( address => OTP) otpList;
    address owner;

    /* CONSTRUCTOR */

    constructor() public payable{
        // Set the owner of the company
        owner = msg.sender;

        // Add it to the admin list
        userList[owner].isNull = false;
        userList[owner].isAdmin = true;
    }

    function rmUser(address _addr) public isAdmin {
        userList[_addr].isNull = false;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
    }

    function addUser(address _addr) public isAdmin {
        userList[_addr].isNull = true;
        userList[_addr].isAdmin = false;
        userList[_addr].isLoggedIn = false;
    }



    function addAdmin(address _addr) public isAdmin {
        // Check that the user is added
        require(userList[_addr].isNull == false,"User does not exist.");
        // We update the user's profile with admin status
        userList[_addr].isAdmin = true;
        // We notify in the blockchain who did it
        emit createAdmin(_addr, msg.sender, now);
    }

}