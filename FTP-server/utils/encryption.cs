using System.Security.Cryptography;
using System.Text;

namespace FTP_server.utils
{
    public static class AES256
    {
        public static string DecryptString(string cipherText, string keyString, string ivString)
        {
            // 将字符串表示的密钥和IV转换为字节数组
            byte[] key = Encoding.UTF8.GetBytes(keyString);
            byte[] iv = Encoding.UTF8.GetBytes(ivString);

            // 创建一个新的Aes对象来表示AES算法
            using (Aes aesAlg = Aes.Create())
            {
                aesAlg.Key = key;
                aesAlg.IV = iv;
                aesAlg.Mode = CipherMode.CBC;
                aesAlg.Padding = PaddingMode.PKCS7;

                // 创建一个解密器来执行流转换
                ICryptoTransform decryptor = aesAlg.CreateDecryptor(aesAlg.Key, aesAlg.IV);

                // 创建必要的流
                using (MemoryStream msDecrypt = new MemoryStream(Convert.FromBase64String(cipherText)))
                {
                    using (CryptoStream csDecrypt = new CryptoStream(msDecrypt, decryptor, CryptoStreamMode.Read))
                    {
                        using (StreamReader srDecrypt = new StreamReader(csDecrypt))
                        {
                            // 从解密流中读取解密后的字节并将它们转换为字符串
                            return srDecrypt.ReadToEnd();
                        }
                    }
                }
            }
        }
    }
}
